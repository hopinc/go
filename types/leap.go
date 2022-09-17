package types

import "errors"

// LeapConnectionState is used to define the state of the Leap connection.
type LeapConnectionState string

const (
	// LeapConnectionStateIdle is the state when the connection is idle.
	LeapConnectionStateIdle LeapConnectionState = "idle"

	// LeapConnectionStateConnecting is the state when the connection is connecting.
	LeapConnectionStateConnecting LeapConnectionState = "connecting"

	// LeapConnectionStateAuthenticating is the state when the connection is authenticating.
	LeapConnectionStateAuthenticating LeapConnectionState = "authenticating"

	// LeapConnectionStateConnected is the state when the connection is connected.
	LeapConnectionStateConnected LeapConnectionState = "connected"

	// LeapConnectionStateErrored is the state when the connection is errored.
	LeapConnectionStateErrored LeapConnectionState = "errored"
)

// LeapScope is used to define the Leap connection scope.
type LeapScope string

const (
	// ScopeProject is the scope for project connections.
	ScopeProject LeapScope = "project"

	// ScopeToken is the scope for token connections.
	ScopeToken LeapScope = "token"
)

// LeapDispatchEventDetails is used to easily pass through some of the data from the dispatch event.
type LeapDispatchEventDetails struct {
	// ChannelID is the ID of the channel.
	ChannelID string `json:"c"`

	// Unicast is if the dispatch event was dispatched via unicast.
	Unicast bool `json:"u"`
}

// LeapInitEvent is used to define the Leap init event.
type LeapInitEvent struct {
	LeapDispatchEventDetails `json:"-"`

	// ConnectionID is the ID of the connection.
	ConnectionID string `json:"cid"`

	// Metadata is the metadata of the connection if applicable.
	Metadata map[string]any `json:"metadata"`

	// Scope is the scope of the connection.
	Scope LeapScope `json:"scope"`

	// Channels is the channels the connection is subscribed to.
	Channels []*ChannelPartial `json:"channels"`
}

// LeapAvailableEvent is used to define the event when a channel is available.
type LeapAvailableEvent struct {
	LeapDispatchEventDetails `json:"-"`

	// Channel is the channel that is available.
	Channel *ChannelPartial `json:"channel"`
}

// LeapUnavailableEvent is used to define the event when a channel is unavailable.
type LeapUnavailableEvent struct {
	LeapDispatchEventDetails `json:"-"`

	// Graceful is if the channel was gracefully unavailable.
	Graceful bool `json:"graceful"`

	// ErrorCode is the error code of the channel.
	ErrorCode string `json:"error_code"`
}

// LeapMessageEvent is used to define the Leap message event. When this is sent, if this is a direct message, ChannelID will
// be blank.
type LeapMessageEvent struct {
	LeapDispatchEventDetails `json:"-"`

	// Data is the user provided event data for the message.
	Data map[string]any `json:"d"`

	// EventName is the name of the event.
	EventName string `json:"e"`
}

// IsDirectMessage returns if this is a direct message.
func (e LeapMessageEvent) IsDirectMessage() bool { return e.ChannelID == "" }

// ExpectedHello is thrown if the first packet after connection is not a hello.
var ExpectedHello = errors.New("expected hello packet after connection")

// LeapAuthorizationError is thrown if the authorization fails.
type LeapAuthorizationError struct {
	Data string
}

// Error returns the error message.
func (e LeapAuthorizationError) Error() string { return e.Data }

// LeapChannelStateUpdateEvent is used to define the channel state update event.
type LeapChannelStateUpdateEvent struct {
	LeapDispatchEventDetails `json:"-"`

	// State is the state of the channel.
	State map[string]any `json:"state"`
}

// LeapStateInfo is the information about the state of the connection.
type LeapStateInfo struct {
	// ConnectionState is the string representation of the connection state.
	ConnectionState LeapConnectionState

	// Err is set if the connection state is errored to define the error that triggered this.
	Err error

	// WillReconnect is set if the connection state is errored to define if a reconnection will be attempted.
	WillReconnect bool
}
