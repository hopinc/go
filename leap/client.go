package leap

import (
	"compress/zlib"
	"encoding/json"
	"errors"
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

// Client is used to define a Leap client. Please use NewClient to create a new client.
type Client struct {
	projectId string
	token     string

	logger Logger

	ws      webSocketImpl
	wsLock  sync.RWMutex
	wsMaker func(string) (webSocketImpl, error)
	url     string

	// Each side of the connection is not thread safe for its own way. Read is only read from one function,
	// but we need a solution here. This is that.
	writeLock sync.Mutex

	state     rwLocker[types.LeapStateInfo]
	initEvent rwLocker[*types.LeapInitEvent]

	channelWaiter eventWaiter[*types.ChannelPartial]

	channelQueue     []*queueDispatcher[types.LeapChannelEvent]
	channelQueueLock sync.RWMutex

	messageQueue     []*queueDispatcher[types.LeapMessageEvent]
	messageQueueLock sync.RWMutex
}

// MessageEventChannel is used to get a channel that will receive all message events.
func (c *Client) MessageEventChannel() <-chan types.LeapMessageEvent {
	ch := make(chan types.LeapMessageEvent)
	c.messageQueueLock.Lock()
	c.messageQueue = append(c.messageQueue, newQueueDispatcher(ch))
	c.messageQueueLock.Unlock()
	return ch
}

// ChannelEventChannel is used to get a channel that will receive all channel events. The only exception to this is events
// that are a reply to functions in this package. Those will be returned from the function.
func (c *Client) ChannelEventChannel() <-chan types.LeapChannelEvent {
	ch := make(chan types.LeapChannelEvent)
	c.channelQueueLock.Lock()
	c.channelQueue = append(c.channelQueue, newQueueDispatcher(ch))
	c.channelQueueLock.Unlock()
	return ch
}

// Closes all queues in the client.
func (c *Client) closeAllQueues() {
	c.channelQueueLock.Lock()
	q := c.channelQueue
	c.channelQueue = nil
	c.channelQueueLock.Unlock()
	for _, v := range q {
		close(v.channel)
	}

	c.messageQueueLock.Lock()
	x := c.messageQueue
	c.messageQueue = nil
	c.messageQueueLock.Unlock()
	for _, v := range x {
		close(v.channel)
	}
}

// Close closes the connection and all channels for events.
func (c *Client) Close() error {
	// Ensure that the websocket is unusable whilst it is being shut down.
	c.wsLock.Lock()
	defer c.wsLock.Unlock()

	// Close the websocket. Keep this error for later.
	err := c.ws.Close()

	// Send any clients waiting for channel relating events that the socket is closed.
	c.channelWaiter.close(net.ErrClosed)

	// Destroy every queue item.
	c.closeAllQueues()

	// Return the error from closing the websocket.
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
func (c *Client) readPayload(heartbeatWriteDuration time.Duration) (*payload, error) {
	err := c.ws.SetReadDeadline(time.Now().Add(heartbeatWriteDuration))
	if err != nil {
		return nil, err
	}
	t, r, err := c.ws.NextReader()
	if err != nil {
		return nil, err
	}
	var p payload
	if t == websocket.BinaryMessage {
		r, err = zlib.NewReader(r)
		if err == nil {
			c.logger.Debug("binary message read - using zlib", nil)
		} else {
			c.logger.Error("binary message read - unable to use zlib", err, nil)
			return nil, err
		}
	} else {
		c.logger.Debug("text message read - not using zlib", nil)
	}
	err = json.NewDecoder(r).Decode(&p)
	if err != nil {
		c.logger.Error("unable to decode json", err, nil)
		return nil, err
	}
	return &p, nil
}

func (c *Client) writePayload(ws webSocketImpl, p *payload) error {
	// Ensure we are only doing 1 write at a time.
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

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

func (c *Client) handleWsError(code int, text string, err error) {
	// Make sure the websocket is killed and set to nil.
	c.wsLock.Lock()
	_ = c.ws.Close()
	c.ws = nil

	// Turn a code 4001 into a AuthorizationError.
	if code == 4001 {
		err = types.LeapAuthorizationError{Data: text}
	}

	// Set the state to error.
	stateInfo := types.LeapStateInfo{
		ConnectionState: types.LeapConnectionStateErrored,
		Err:             err,
		WillReconnect:   code != 4001,
	}
	c.state.set(stateInfo)

	// Error log the information.
	c.logger.Error("connection failure", err, map[string]any{
		"will_reconnect": stateInfo.WillReconnect,
	})

	// Make sure all channel waiters are closed.
	c.channelWaiter.close(err)

	// Unlock the websocket. We will relock it in just a second.
	c.wsLock.Unlock()

	// Check if this is a close error.
	if code == 4006 {
		// Change the url before reconnect.
		c.url = text
	}
	if code == 4001 {
		// If the code is 4001, this means that the connection was closed on purpose. This is not something we should
		// reconnect for.
		c.closeAllQueues()
	} else {
		// Attempt looping until we reconnect.
		for {
			c.logger.Warn("connection dropped - attempting to reconnect", nil)
			err = c.connect(true)
			if err == nil {
				// We are ready to rumble!
				c.logger.Info("reconnected", nil)
				return
			}
			c.logger.Error("failed to reconnect - backing off for a second", err, nil)
			time.Sleep(time.Second)
		}
	}
}

// This is used to define an event for channel messages.
type dispatchEvent struct {
	types.LeapDispatchEventDetails `json:",inline"`

	// DispatchEventCode is the code of the dispatch event.
	DispatchEventCode string `json:"e"`

	// Data is the data of the dispatch event.
	Data json.RawMessage `json:"d"`
}

// Unmarshals the data into the given interface.
func unmarshalPacket(e dispatchEvent, x any) error {
	err := json.Unmarshal(e.Data, x)
	if err != nil {
		return err
	}
	reflect.Indirect(reflect.ValueOf(x)).FieldByName("LeapDispatchEventDetails").
		Set(reflect.ValueOf(e.LeapDispatchEventDetails))
	return nil
}

// InitEvent is used to return the init event. Can be nil if it is not sent.
func (c *Client) InitEvent() *types.LeapInitEvent {
	return c.initEvent.get()
}

// Used to handle dispatching events.
func (c *Client) dispatchEvent(r json.RawMessage) {
	var x dispatchEvent
	err := json.Unmarshal(r, &x)
	if err != nil {
		return
	}

	c.logger.Debug("dispatch event", map[string]any{
		"event_code": x.DispatchEventCode,
		"channel_id": x.ChannelID,
		"unicast":    x.Unicast,
		"data":       x.Data,
	})

	switch x.DispatchEventCode {
	case "INIT":
		var e types.LeapInitEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		c.initEvent.set(&e)
		c.state.set(types.LeapStateInfo{ConnectionState: types.LeapConnectionStateConnected})
		c.logger.Info("init event received - we are connected", nil)
	case "AVAILABLE":
		var e types.LeapAvailableEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		if ok := c.channelWaiter.signal(e.Channel.ID, e.Channel, nil); !ok {
			c.channelQueueLock.RLock()
			for _, v := range c.channelQueue {
				v.dispatch(e)
			}
			c.channelQueueLock.RUnlock()
		}
	case "UNAVAILABLE":
		var e types.LeapUnavailableEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		if ok := c.channelWaiter.signal(e.ChannelID, nil, e); !ok {
			c.channelQueueLock.RLock()
			for _, v := range c.channelQueue {
				v.dispatch(e)
			}
			c.channelQueueLock.RUnlock()
		}
	case "MESSAGE", "DIRECT_MESSAGE": // MESSAGE and DIRECT_MESSAGE are the same inside packet.
		var e types.LeapMessageEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		c.messageQueueLock.RLock()
		for _, v := range c.messageQueue {
			v.dispatch(e)
		}
		c.messageQueueLock.RUnlock()
	case "STATE_UPDATE":
		var e types.LeapChannelStateUpdateEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		c.channelQueueLock.RLock()
		for _, v := range c.channelQueue {
			v.dispatch(e)
		}
		c.channelQueueLock.RUnlock()
	case "PIPE_ROOM_AVAILABLE":
		var e types.LeapPipeRoomAvailableEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		c.channelQueueLock.RLock()
		for _, v := range c.channelQueue {
			v.dispatch(e)
		}
		c.channelQueueLock.RUnlock()
	case "PIPE_ROOM_UPDATE":
		var e types.LeapPipeRoomUpdateEvent
		if err = unmarshalPacket(x, &e); err != nil {
			return
		}
		c.channelQueueLock.RLock()
		for _, v := range c.channelQueue {
			v.dispatch(e)
		}
		c.channelQueueLock.RUnlock()
	default:
		c.logger.Warn("unknown dispatch event", map[string]any{
			"event_code": x.DispatchEventCode,
		})
	}
}

// Subscribe subscribes to a channel.
func (c *Client) Subscribe(channelId string) (*types.ChannelPartial, error) {
	c.wsLock.RLock()
	ws := c.ws
	c.wsLock.RUnlock()
	if ws == nil {
		return nil, net.ErrClosed
	}
	err := c.writePayload(ws, &payload{
		Op: 0,
		Data: rawify(dispatchEvent{
			LeapDispatchEventDetails: types.LeapDispatchEventDetails{
				ChannelID: channelId,
				Unicast:   false,
			},
			DispatchEventCode: "SUBSCRIBE",
		}),
	})
	if err != nil {
		return nil, err
	}
	return c.channelWaiter.wait(channelId)
}

// Defines the read loop.
func (c *Client) readLoop(ws webSocketImpl, d time.Duration) {
	for {
		// Read the payload.
		p, err := c.readPayload(d)
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				// Call the close handler and then return.
				c.handleWsError(closeErr.Code, closeErr.Text, err)
				return
			}

			if errors.Is(err, net.ErrClosed) {
				// Give up on this connection.
				return
			}

			if _, ok := err.(net.Error); ok {
				// Call the close handler with -1 and then return. This technically is not a close, but we want to
				// handle it the same way.
				c.handleWsError(-1, "", err)
				return
			}

			// Loop around again.
			c.logger.Error("unknown error reading payload - looping around", err, nil)
			continue
		}

		// Handle the payload.
		switch p.Op {
		case 0:
			// Dispatch the event.
			c.logger.Debug("opcode 0 received", map[string]any{
				"data": p.Data,
			})
			c.dispatchEvent(p.Data)
		case 3:
			// Reply with a heartbeat.
			c.logger.Debug("heartbeat received", nil)
			go func() {
				_ = c.writePayload(ws, &payload{
					Op:   3,
					Data: p.Data,
				})
				c.logger.Debug("heartbeat reply - sent", nil)
			}()
		case 4:
			c.logger.Debug("heartbeat ack received", nil)
		default:
			// Unknown op code.
			c.logger.Warn("unknown packet received", map[string]any{
				"op":   p.Op,
				"data": p.Data,
			})
		}
	}
}

// Defines the heartbeat loop.
func (c *Client) heartbeatLoop(ws webSocketImpl, interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		for {
			c.logger.Debug("heartbeat ticker - loop", nil)
			go func() {
				err := c.writePayload(ws, &payload{
					Op:   3,
					Data: rawify(map[string]string{"tag": ""}),
				})
				if err == nil {
					// Log that it was sent.
					c.logger.Debug("heartbeat ticker - sent", nil)
				} else {
					// Log that it failed.
					c.logger.Error("heartbeat ticker - failed", err, nil)

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

func (c *Client) connect(reconnect bool) error {
	// Take the websocket mutex.
	c.wsLock.Lock()
	defer c.wsLock.Unlock()

	// Check if we are already connected.
	if c.ws != nil {
		return nil
	}

	// Set the state to connecting.
	c.state.set(types.LeapStateInfo{ConnectionState: types.LeapConnectionStateConnecting})

	// Make a new websocket.
	var err error
	c.ws, err = c.wsMaker(c.url)
	if err != nil {
		_ = c.ws.Close()
		c.ws = nil
		c.state.set(types.LeapStateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: err, WillReconnect: reconnect})
		return err
	}

	// Read the first payload.
	p, err := c.readPayload(time.Second * 10)
	if err != nil {
		// Unable to recover from whatever happened in the read event.
		_ = c.ws.Close()
		c.ws = nil
		c.state.set(types.LeapStateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: err, WillReconnect: reconnect})
		return err
	}

	// Validate this is a hello message.
	if p.Op != 1 {
		_ = c.ws.Close()
		c.ws = nil
		c.state.set(types.LeapStateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: types.ExpectedHello, WillReconnect: reconnect})
		return types.ExpectedHello
	}
	type hello struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	}
	var h hello
	if err = json.Unmarshal(p.Data, &h); err != nil {
		c.state.set(types.LeapStateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: err, WillReconnect: reconnect})
		_ = c.ws.Close()
		c.ws = nil
		return err
	}

	// Send the identify payload.
	c.state.set(types.LeapStateInfo{ConnectionState: types.LeapConnectionStateAuthenticating})
	err = c.writePayload(c.ws, &payload{
		Op: 2,
		Data: rawify(map[string]string{
			"token":      c.token,
			"project_id": c.projectId,
		}),
	})
	if err != nil {
		c.state.set(types.LeapStateInfo{ConnectionState: types.LeapConnectionStateErrored, Err: err, WillReconnect: reconnect})
		_ = c.ws.Close()
		c.ws = nil
		return err
	}

	// Start the reading loop.
	go c.readLoop(c.ws, (time.Millisecond*time.Duration(h.HeartbeatInterval))+(time.Second*5))

	// Start the heartbeat loop.
	c.heartbeatLoop(c.ws, h.HeartbeatInterval)

	// Return no errors.
	return nil
}

// Connect is used to connect to the Leap server.
func (c *Client) Connect() error {
	return c.connect(false)
}

// State returns the state of the websocket.
func (c *Client) State() types.LeapStateInfo {
	return c.state.get()
}

// AddStateUpdateListener adds a handler to be called when the state changes.
func (c *Client) AddStateUpdateListener(handler func(types.LeapStateInfo)) {
	c.state.addListener(handler)
}

// NewClient is used to create a new client. If the specified logger is nil, the library will not log any data.
func NewClient(projectId, token string, l Logger) *Client {
	if l == nil {
		l = NopLogger{}
	}
	return &Client{
		projectId: projectId,
		token:     token,
		logger:    l,
		state: rwLocker[types.LeapStateInfo]{unsafeValue: types.LeapStateInfo{
			ConnectionState: types.LeapConnectionStateIdle,
		}},
		wsMaker: newWebSocketImpl,
		url:     "wss://leap.hop.io/ws?encoding=json&compression=zlib",
	}
}
