package types

// ChannelType is used to define the type of the Hop channel.
type ChannelType string

const (
	// ChannelTypePrivate is used to define a private Hop channel.
	ChannelTypePrivate ChannelType = "private"

	// ChannelTypePublic is used to define a public Hop channel.
	ChannelTypePublic ChannelType = "public"

	// ChannelTypeUnprotected is used to define an unprotected Hop channel.
	ChannelTypeUnprotected ChannelType = "unprotected"
)

// ChannelPartial is used to define a partial channel. Whilst currently not used by the hop-go library,
// this may be desirable for things importing the types package such as a leap client.
type ChannelPartial struct {
	// ID is the ID of the channel.
	ID string `json:"id"`

	// State is any state metadata associated with the channel.
	State map[string]any `json:"state"`

	// Type is the type of this channel.
	Type ChannelType `json:"type"`
}

// Channel is used to define the main structure of a channel.
type Channel struct {
	ChannelPartial `json:",inline"`

	// Project is the project it is associated with.
	Project *Project `json:"project"`

	// Capabilities is the capabilities of this channel.
	Capabilities int `json:"capabilities"`

	// CreatedAt is when this channel was created.
	CreatedAt Timestamp `json:"created_at"`
}

// Stats is used to define the stats for a channel.
type Stats struct {
	OnlineCount int `json:"online_count"`
}

// ChannelToken is used to define the token for a channel.
type ChannelToken struct {
	// ID is the ID of the token.
	ID string `json:"id"`

	// State is any state metadata associated with the token.
	State map[string]any `json:"state"`

	// ProjectID is the project ID associated with the token.
	ProjectID string `json:"project_id"`

	// IsOnline is whether the token is online (e.g.: active heartbeat and connected to leap).
	IsOnline bool `json:"is_online"`
}
