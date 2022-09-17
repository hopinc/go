package leap

import (
	"compress/zlib"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hopinc/hop-go/types"
)

type webSocketImpl interface {
	NextReader() (int, io.Reader, error)
	NextWriter(int) (io.WriteCloser, error)
	SetReadDeadline(time.Time) error
	Close() error
}

func newWebSocketImpl(url string) (webSocketImpl, error) {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// TODO: support logging why it does things

type rwLocker[T any] struct {
	unsafeValue T
	mu          sync.RWMutex

	// Setting is the perfect time to fire events. This is because a mutex has to be locked anyway.
	changes []func(T)
}

func (r *rwLocker[T]) get() T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.unsafeValue
}

func (r *rwLocker[T]) set(value T) {
	r.mu.Lock()
	r.unsafeValue = value
	changes := r.changes
	r.mu.Unlock()
	for _, f := range changes {
		go f(value)
	}
}

func (r *rwLocker[T]) addListener(f func(T)) {
	r.mu.Lock()
	r.changes = append(r.changes, f)
	r.mu.Unlock()
}

// StateInfo is the information about the state of the connection.
type StateInfo struct {
	// ConnectionState is the string representation of the connection state.
	ConnectionState types.LeapConnectionState

	// Err is set if the connection state is errored to define the error that triggered this.
	Err error

	// WillReconnect is set if the connection state is errored to define if a reconnection will be attempted.
	WillReconnect bool
}

// Client is used to define a Leap client. Please use NewClient to create a new client.
type Client struct {
	projectId string
	token     string

	ws      webSocketImpl
	wsLock  sync.RWMutex
	wsMaker func(string) (webSocketImpl, error)
	url     string

	// Each side of the connection is not thread safe for its own way. Read is only read from one function,
	// but we need a solution here. This is that.
	writeLock sync.Mutex

	state     rwLocker[StateInfo]
	initEvent rwLocker[*InitEvent]

	channelWaiter eventWaiter[*types.ChannelPartial]

	messageQueue     []*queueDispatcher[MessageEvent]
	messageQueueLock sync.RWMutex
}

// MessageEventChannel is used to get a channel that will receive all message events.
func (w *Client) MessageEventChannel() <-chan MessageEvent {
	ch := make(chan MessageEvent)
	w.messageQueueLock.Lock()
	w.messageQueue = append(w.messageQueue, newQueueDispatcher(ch))
	w.messageQueueLock.Unlock()
	return ch
}

// Close closes the connection and all channels for events.
func (w *Client) Close() error {
	w.wsLock.Lock()
	defer w.wsLock.Unlock()
	err := w.ws.Close()
	w.channelWaiter.close(net.ErrClosed)
	w.messageQueueLock.Lock()
	q := w.messageQueue
	w.messageQueueLock.Unlock()
	for _, v := range q {
		close(v.channel)
	}
	return err
}

type payload struct {
	Op   int             `json:"op"`
	Data json.RawMessage `json:"d"`
}

func rawify(data any) json.RawMessage {
	j, _ := json.Marshal(data)
	return j
}

// If payload and error is both nil, this means there was an error, but it was handled. Just call this again (unless
// it is in connect, then return)!
func (w *Client) readPayload(heartbeatWriteDuration time.Duration) (*payload, error) {
	err := w.ws.SetReadDeadline(time.Now().Add(heartbeatWriteDuration))
	if err != nil {
		return nil, err
	}
	t, r, err := w.ws.NextReader()
	if err != nil {
		return nil, err
	}
	var p payload
	if t == websocket.BinaryMessage {
		r, err = zlib.NewReader(r)
		if err != nil {
			return nil, err
		}
	}
	err = json.NewDecoder(r).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (w *Client) writePayload(ws webSocketImpl, p *payload) error {
	// Ensure we are only doing 1 write at a time.
	w.writeLock.Lock()
	defer w.writeLock.Unlock()

	// Get the next writer.
	wr, err := ws.NextWriter(websocket.TextMessage)
	if err != nil {
		// ...or not.
		return err
	}
	defer wr.Close()

	// Make a zlib writer and then write json to it.
	err = json.NewEncoder(wr).Encode(p)
	if err != nil {
		return err
	}

	// Tell the error channel everything was dandy.
	return nil
}

func (w *Client) handleWsError(code int, text string, err error) {
	// Make sure the websocket is killed and set to nil.
	w.wsLock.Lock()
	_ = w.ws.Close()
	w.ws = nil

	// Turn a code 4001 into a AuthorizationError.
	if code == 4001 {
		err = AuthorizationError{data: text}
	}

	// Set the state to error.
	w.state.set(StateInfo{
		ConnectionState: types.LeapConnectionStateErrored,
		Err:             err,
		WillReconnect:   code != 4001,
	})

	// Make sure all channel waiters are closed.
	w.channelWaiter.close(err)

	// Unlock the websocket. We will relock it in just a second.
	w.wsLock.Unlock()

	// Check if this is a close error.
	if code == 4006 {
		// Change the url before reconnect.
		w.url = text
	}
	if code == 4001 {
		// If the code is 4001, this means that the connection was closed on purpose. This is not something we should
		// reconnect for.
		w.messageQueueLock.Lock()
		q := w.messageQueue
		w.messageQueueLock.Unlock()
		for _, v := range q {
			close(v.channel)
		}
	} else {
		// Attempt looping until we reconnect.
		for {
			err = w.connect(true)
			if err == nil {
				// We are ready to rumble!
				return
			}
			time.Sleep(time.Second)
		}
	}
}

// Unmarshals the data into the given interface.
func unmarshalPacket(e dispatchEvent, x any) error {
	err := json.Unmarshal(e.Data, x)
	if err != nil {
		return err
	}
	reflect.Indirect(reflect.ValueOf(x)).FieldByName("DispatchEventDetails").
		Set(reflect.ValueOf(e.DispatchEventDetails))
	return nil
}

// InitEvent is used to return the init event. Can be nil if it is not sent.
func (w *Client) InitEvent() *InitEvent {
	return w.initEvent.get()
}

// Used to handle dispatching events.
func (w *Client) dispatchEvent(r json.RawMessage) {
	var x dispatchEvent
	err := json.Unmarshal(r, &x)
	if err != nil {
		return
	}

	// TODO: make a nicer channel events system for if a channel randomly becomes available/unavailable.
	switch x.DispatchEventCode {
	case "INIT":
		var e InitEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		w.initEvent.set(&e)
		w.state.set(StateInfo{ConnectionState: types.LeapConnectionStateConnected})
	case "AVAILABLE":
		var e AvailableEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		w.channelWaiter.signal(e.Channel.ID, e.Channel, nil)
	case "UNAVAILABLE":
		var e UnavailableEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		// TODO: Make this a better error.
		w.channelWaiter.signal(e.ChannelID, nil, fmt.Errorf("%v", e))
	case "MESSAGE", "DIRECT_MESSAGE": // MESSAGE and DIRECT_MESSAGE are the same packet.
		var e MessageEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		w.messageQueueLock.RLock()
		for _, v := range w.messageQueue {
			v.dispatch(e)
		}
		w.messageQueueLock.RUnlock()
	case "STATE_UPDATE":
		var e ChannelStateUpdateEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		fmt.Println(e)
	}
}

// Subscribe subscribes to a channel.
func (w *Client) Subscribe(channelId string) (*types.ChannelPartial, error) {
	w.wsLock.RLock()
	ws := w.ws
	w.wsLock.RUnlock()
	if ws == nil {
		return nil, net.ErrClosed
	}
	err := w.writePayload(ws, &payload{
		Op: 0,
		Data: rawify(dispatchEvent{
			DispatchEventDetails: DispatchEventDetails{
				ChannelID: channelId,
				Unicast:   false,
			},
			DispatchEventCode: "SUBSCRIBE",
		}),
	})
	if err != nil {
		return nil, err
	}
	return w.channelWaiter.wait(channelId)
}

// Defines the read loop.
func (w *Client) readLoop(ws webSocketImpl, d time.Duration) {
	for {
		// Read the payload.
		p, err := w.readPayload(d)
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				// Call the close handler and then return.
				w.handleWsError(closeErr.Code, closeErr.Text, err)
				return
			}

			if errors.Is(err, net.ErrClosed) {
				// Give up on this connection.
				return
			}

			if _, ok := err.(net.Error); ok {
				// Call the close handler with -1 and then return. This technically is not a close, but we want to
				// handle it the same way.
				w.handleWsError(-1, "", err)
				return
			}

			// Loop around again.
			continue
		}

		// Handle the payload.
		switch p.Op {
		case 0:
			// Dispatch the event.
			w.dispatchEvent(p.Data)
		case 3:
			// Reply with a heartbeat.
			_ = w.writePayload(ws, &payload{
				Op:   3,
				Data: p.Data,
			})
		}
	}
}

// Defines the heartbeat loop.
func (w *Client) heartbeatLoop(ws webSocketImpl, interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		for {
			go func() {
				err := w.writePayload(ws, &payload{
					Op:   3,
					Data: rawify(map[string]string{"tag": ""}),
				})
				if err != nil {
					// Stop this heart beating. The read loop will take over error handling.
					t.Stop()
				}
			}()
			_, ok := <-t.C
			if !ok {
				return
			}
		}
	}()
}

func (w *Client) connect(reconnect bool) error {
	// Take the websocket mutex.
	w.wsLock.Lock()
	defer w.wsLock.Unlock()

	// Check if we are already connected.
	if w.ws != nil {
		return nil
	}

	// Set the state to connecting.
	w.state.set(StateInfo{ConnectionState: types.LeapConnectionStateConnecting})

	// Make a new websocket.
	var err error
	w.ws, err = w.wsMaker(w.url)
	if err != nil {
		w.state.set(StateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: err, WillReconnect: reconnect})
		return err
	}

	// Read the first payload.
	p, err := w.readPayload(time.Second * 10)
	if err != nil {
		// Unable to recover from whatever happened in the read event.
		w.state.set(StateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: err, WillReconnect: reconnect})
		return err
	}

	// Validate this is a hello message.
	if p.Op != 1 {
		_ = w.ws.Close()
		w.ws = nil
		w.state.set(StateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: ExpectedHello, WillReconnect: reconnect})
		return ExpectedHello
	}
	type hello struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	}
	var h hello
	if err = json.Unmarshal(p.Data, &h); err != nil {
		w.state.set(StateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: err, WillReconnect: reconnect})
		_ = w.ws.Close()
		w.ws = nil
		return err
	}

	// Send the identify payload.
	w.state.set(StateInfo{ConnectionState: types.LeapConnectionStateAuthenticating})
	err = w.writePayload(w.ws, &payload{
		Op: 2,
		Data: rawify(map[string]string{
			"token":      w.token,
			"project_id": w.projectId,
		}),
	})
	if err != nil {
		w.state.set(StateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: ExpectedHello, WillReconnect: reconnect})
		return err
	}

	// Start the reading loop.
	go w.readLoop(w.ws, (time.Millisecond*time.Duration(h.HeartbeatInterval))+(time.Second*5))

	// Start the heartbeat loop.
	w.heartbeatLoop(w.ws, h.HeartbeatInterval)

	// Return no errors.
	return nil
}

// Connect is used to connect to the Leap server.
func (w *Client) Connect() error {
	return w.connect(false)
}

// State returns the state of the websocket.
func (w *Client) State() StateInfo {
	return w.state.get()
}

// AddStateUpdateListener adds a handler to be called when the state changes.
func (w *Client) AddStateUpdateListener(handler func(StateInfo)) {
	w.state.addListener(handler)
}

// NewClient is used to create a new client.
func NewClient(projectId, token string) *Client {
	return &Client{
		projectId: projectId,
		token:     token,
		state:     rwLocker[StateInfo]{unsafeValue: StateInfo{ConnectionState: types.LeapConnectionStateIdle}},
		wsMaker:   newWebSocketImpl,
		url:       "wss://leap.hop.io/ws?encoding=json&compression=zlib",
	}
}
