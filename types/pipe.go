package types

// IngestProtocol is the type used for supported Pipe ingest protocols.
type IngestProtocol string

const (
	// IngestProtocolRTMP is used to define the RTMP ingest protocol.
	IngestProtocolRTMP IngestProtocol = "rtmp"
)

// DeliveryProtocol is the type used for supported Pipe delivery protocols.
type DeliveryProtocol string

const (
	// DeliveryProtocolHLS is used to define the HLS delivery protocol.
	DeliveryProtocolHLS DeliveryProtocol = "hls"

	// DeliveryProtocolWebRTC is used to define the WebRTC delivery protocol.
	DeliveryProtocolWebRTC DeliveryProtocol = "webrtc"
)

// RoomState is used to define the state of a room.
type RoomState string

const (
	// RoomStateLive is used to define a live room.
	RoomStateLive RoomState = "live"

	// RoomStateOffline is used to define an offline room.
	RoomStateOffline RoomState = "offline"
)

// Room is used to define the main structure of a Pipe room.
type Room struct {
	// ID is the ID of the room.
	ID string `json:"id"`

	// Name is the name of the room.
	Name string `json:"name"`

	// CreatedAt is when this room was created.
	CreatedAt Timestamp `json:"created_at"`

	// IngestProtocol is the protocol you can stream with.
	IngestProtocol IngestProtocol `json:"ingest_protocol"`

	// DeliveryProtocols are the protocols that are supported by this room to the client.
	DeliveryProtocols []DeliveryProtocol `json:"delivery_protocols"`

	// JoinToken is the token to subscribe to this room.
	JoinToken string `json:"join_token"`

	// IngestRegion is the region that the stream URL is located in.
	IngestRegion Region `json:"ingest_region"`

	// State is the current state of the room.
	State RoomState `json:"state"`
}

// HLSConfig is used to define the HLS configuration for a room.
type HLSConfig struct {
	WCLDelay                int    `json:"wcl_delay"`
	ArtificialDelay         int    `json:"artificial_delay"`
	MaxPlayoutBitratePreset string `json:"max_playout_bitrate_preset"`
}

// RoomCreationOptions is used to define the options for creating a room.
type RoomCreationOptions struct {
	// ProjectID is the ID of the project that this gateway is for. Can be blank if using a project token.
	ProjectID string `json:"-"`

	// Name is the name of the room.
	Name string `json:"name"`

	// DeliveryProtocols are the protocols that are supported by this room to the client.
	DeliveryProtocols []DeliveryProtocol `json:"delivery_protocols"`

	// Ephemeral defines whether the room is ephemeral or not.
	Ephemeral bool `json:"ephemeral"`

	// Region is used to define the room region.
	Region Region `json:"region"`

	// IngestProtocol is the protocol you can stream with.
	IngestProtocol IngestProtocol `json:"ingest_protocol"`

	// HLSConfig is the configuration for HLS delivery. This can be nil.
	HLSConfig HLSConfig `json:"llhls_config"`
}
