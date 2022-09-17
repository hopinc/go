package types

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
