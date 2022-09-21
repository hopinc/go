package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// A very crude websocket implementation for the demo. It is just meant to be real simple.

type wsWrapper struct {
	ws *websocket.Conn

	writeLock sync.Mutex

	onClose   func()
	onCloseMu sync.Mutex
}

func (w *wsWrapper) readLoop() {
	for {
		err := w.ws.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
		if err != nil {
			w.onCloseMu.Lock()
			x := w.onClose
			w.onCloseMu.Unlock()
			if x != nil {
				x()
			}
			return
		}
		_, _, err = w.ws.ReadMessage()
		if err != nil {
			w.onCloseMu.Lock()
			x := w.onClose
			w.onCloseMu.Unlock()
			if x != nil {
				x()
			}
			return
		}
	}
}

func (w *wsWrapper) writePing() error {
	w.writeLock.Lock()
	defer w.writeLock.Unlock()
	return w.ws.WriteMessage(websocket.BinaryMessage, []byte{0})
}

func (w *wsWrapper) writeData(data []byte) error {
	w.writeLock.Lock()
	defer w.writeLock.Unlock()
	return w.ws.WriteMessage(websocket.BinaryMessage, append([]byte{1}, data...))
}

var upgrader = &websocket.Upgrader{}

func upgradeWs(w http.ResponseWriter, req *http.Request, onClose func()) (*wsWrapper, error) {
	upgrade, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return nil, err
	}
	wr := &wsWrapper{ws: upgrade, onClose: onClose}
	go wr.readLoop()
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			err := wr.writePing()
			if err != nil {
				return
			}
		}
	}()
	return wr, nil
}

var (
	nextWsId           uintptr
	connectedClients   = map[uintptr]*wsWrapper{}
	connectedClientsMu sync.RWMutex
)

func publishMessage(b []byte) {
	connectedClientsMu.RLock()
	clients := make([]*wsWrapper, 0, len(connectedClients))
	for _, client := range connectedClients {
		clients = append(clients, client)
	}
	connectedClientsMu.RUnlock()
	for _, v := range clients {
		go func(v *wsWrapper) {
			_ = v.writeData(b)
		}(v)
	}
}
