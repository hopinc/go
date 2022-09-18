package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/hopinc/hop-go/leap"
	"github.com/hopinc/hop-go/types"
)

func subscribeChannel(w http.ResponseWriter, channel string, c *leap.Client) {
	res, err := c.Subscribe(channel)
	status := 400
	var body map[string]any
	if err == nil {
		body = map[string]any{"channel": res}
	} else {
		body = map[string]any{"error": err.Error()}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func main() {
	// Make sure we were booted with a project ID and leap token.
	if len(os.Args) != 3 {
		panic("Usage: leap-client <project id> <leap token>")
	}
	projectId := os.Args[1]
	leapToken := os.Args[2]

	// Defines the leap client.
	c := leap.NewClient(projectId, leapToken, leap.NewFmtLogger())

	// Start a http server on port 3000.
	go func() {
		err := http.ListenAndServe(":3000", routeHttp(c))
		if err != nil {
			panic(err)
		}
	}()

	// Add event handler for state events.
	c.AddStateUpdateListener(func(info types.LeapStateInfo) {
		if info.ConnectionState == types.LeapConnectionStateErrored && !info.WillReconnect {
			// If we are not going to reconnect, we should exit.
			panic(info.Err)
		}

		// Publish a string for the connection state.
		publishMessage([]byte(`"` + info.ConnectionState + `"`))
	})

	// Start the leap client.
	if err := c.Connect(); err != nil {
		panic(err)
	}

	go func() {
		// Get the channel for channel events. These channels work like any other Go channel in that the reader should not close them,
		// but the client can when it is closed. However, each call to the function does not share the same channel. This
		// means you can make multiple and consume them in multiple places, just do not make one per client or something.
		// The rationale for using channels here is that it makes it clear that it is ordered data.
		ch := c.ChannelEventChannel()

		// Read each channel event.
		for {
			x, ok := <-ch
			if !ok {
				// You should check if the read was okay.
				break
			}

			// LeapChannelEvent is an any type that can be one of LeapUnavailableEvent, LeapAvailableEvent, LeapChannelStateUpdateEvent,
			// LeapPipeRoomAvailableEvent, and LeapPipeRoomUpdateEvent.
			switch x.(type) {
			case types.LeapUnavailableEvent:
			case types.LeapAvailableEvent:
			case types.LeapChannelStateUpdateEvent:
			case types.LeapPipeRoomAvailableEvent:
			case types.LeapPipeRoomUpdateEvent:
			}

			// These events will always JSON marshal.
			b, err := json.Marshal(map[string]any{
				"channel": x,
			})
			if err != nil {
				panic(err)
			}

			// Send the event to the websockets.
			publishMessage(b)
		}
	}()

	// Get a message channel. These channels work like any other Go channel in that the reader should not close them,
	// but the client can when it is closed. However, each call to the function does not share the same channel. This
	// means you can make multiple and consume them in multiple places, just do not make one per client or something.
	// The rationale for using channels here is that it makes it clear that it is ordered data.
	ch := c.MessageEventChannel()

	// Read each message.
	for {
		x, ok := <-ch
		if !ok {
			// You should check if the read was okay.
			break
		}

		// These events will always JSON marshal.
		b, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}

		// Send the event to the websockets.
		publishMessage(b)
	}

	// Sleep for half a second.
	time.Sleep(time.Millisecond * 500)
}
