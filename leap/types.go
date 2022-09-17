package leap

import (
	"encoding/json"
	"errors"

	"github.com/hopinc/hop-go/types"
)

// ConnectionState is used to define the state of the connection.
type ConnectionState string

const (
	// ConnectionStateIdle is the state when the connection is idle.
	ConnectionStateIdle ConnectionState = "idle"

	// ConnectionStateConnecting is the state when the connection is connecting.
	ConnectionStateConnecting ConnectionState = "connecting"

	// ConnectionStateAuthenticating is the state when the connection is authenticating.
	ConnectionStateAuthenticating ConnectionState = "authenticating"

	// ConnectionStateConnected is the state when the connection is connected.
	ConnectionStateConnected ConnectionState = "connected"

	// ConnectionStateErrored is the state when the connection is errored.
	ConnectionStateErrored ConnectionState = "errored"
)

// DispatchEventDetails is used to easily pass through some of the data from the dispatch event.
type DispatchEventDetails struct {
	// ChannelID is the ID of the channel.
	ChannelID string `json:"c"`

	// Unicast is if the dispatch event was dispatched via unicast.
	Unicast bool `json:"u"`
}

// This is used to define an event for channel messages.
type dispatchEvent struct {
	DispatchEventDetails `json:",inline"`

	// DispatchEventCode is the code of the dispatch event.
	DispatchEventCode string `json:"e"`

	// Data is the data of the dispatch event.
	Data json.RawMessage `json:"d"`
}

// Scope is used to define the connection scope.
type Scope string

const (
	// ScopeProject is the scope for project connections.
	ScopeProject Scope = "project"

	// ScopeToken is the scope for token connections.
	ScopeToken Scope = "token"
)

// InitEvent is used to define the init event.
type InitEvent struct {
	DispatchEventDetails `json:"-"`

	// ConnectionID is the ID of the connection.
	ConnectionID string `json:"cid"`

	// Metadata is the metadata of the connection if applicable.
	Metadata map[string]any `json:"metadata"`

	// Scope is the scope of the connection.
	Scope Scope `json:"scope"`

	// Channels is the channels the connection is subscribed to.
	Channels []*types.ChannelPartial `json:"channels"`
}

// AvailableEvent is used to define the event when a channel is available.
type AvailableEvent struct {
	DispatchEventDetails `json:"-"`

	// Channel is the channel that is available.
	Channel *types.ChannelPartial `json:"channel"`
}

// UnavailableEvent is used to define the event when a channel is unavailable.
type UnavailableEvent struct {
	DispatchEventDetails `json:"-"`

	// Graceful is if the channel was gracefully unavailable.
	Graceful bool `json:"graceful"`

	// ErrorCode is the error code of the channel.
	ErrorCode string `json:"error_code"`
}

// MessageEvent is used to define the message event. When this is sent, if this is a direct message, ChannelID will
// be blank.
type MessageEvent struct {
	DispatchEventDetails `json:"-"`

	// Data is the user provided event data for the message.
	Data map[string]any `json:"d"`

	// EventName is the name of the event.
	EventName string `json:"e"`
}

// IsDirectMessage returns if this is a direct message.
func (e MessageEvent) IsDirectMessage() bool { return e.ChannelID == "" }

// ExpectedHello is thrown if the first packet after connection is not a hello.
var ExpectedHello = errors.New("expected hello packet after connection")

// AuthorizationError is thrown if the authorization fails.
type AuthorizationError struct {
	data string
}

// Error returns the error message.
func (e AuthorizationError) Error() string { return e.data }

// ChannelStateUpdateEvent is used to define the channel state update event.
type ChannelStateUpdateEvent struct {
	DispatchEventDetails `json:"-"`

	// State is the state of the channel.
	State map[string]any `json:"state"`
}
