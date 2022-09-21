package main

import (
	_ "embed"
	"github.com/hopinc/hop-go/leap"
	"net/http"
	"sync/atomic"
)

// A very crude http router implementation for the demo. It is just meant to be real simple.

//go:embed index.html
var indexHtml []byte

//go:embed client.js
var clientJs []byte

func routeHttp(c *leap.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write(indexHtml)
		case "/client.js":
			w.Header().Set("Content-Type", "application/javascript")
			_, _ = w.Write(clientJs)
		case "/ws":
			next := atomic.AddUintptr(&nextWsId, 1)
			u, err := upgradeWs(w, req, func() {
				connectedClientsMu.Lock()
				delete(connectedClients, next)
				connectedClientsMu.Unlock()
			})
			if err == nil {
				connectedClientsMu.Lock()
				connectedClients[next] = u
				connectedClientsMu.Unlock()
			}
		case "/subscribe":
			// Get the channel query param.
			channel := req.URL.Query().Get("channel")
			if channel == "" {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("channel query param is required"))
				return
			}
			subscribeChannel(w, channel, c)
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not found"))
		}
	}
}
