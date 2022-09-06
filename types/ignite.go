package types

// GatewayType is the type of the gateway.
type GatewayType string

const (
	// GatewayTypeInternal is used to define a gateway that can only be accessed from within a projects network.
	GatewayTypeInternal GatewayType = "internal"

	// GatewayTypeExternal is used to define a gateway that can be accessed from the internet.
	GatewayTypeExternal GatewayType = "external"
)

// GatewayProtocol is the protocol of the gateway.
type GatewayProtocol string

const (
	// GatewayProtocolHTTP is used to define a gateway that uses the HTTP protocol.
	GatewayProtocolHTTP GatewayProtocol = "http"
)

// DomainState is the state of the domain.
type DomainState string

const (
	// DomainStatePending is used to define a domain that is pending.
	DomainStatePending DomainState = "pending"

	// DomainNameValidCname is used to define a domain that is valid and has a CNAME record.
	DomainNameValidCname DomainState = "valid_cname"

	// DomainNameSSLActive is used to define a domain that has a valid SSL certificate.
	DomainNameSSLActive DomainState = "ssl_active"
)

// Domain is used to define a domain in Ignite.
type Domain struct {
	// ID is the ID of the domain.
	ID string `json:"id"`

	// Domain is the full name of the domain.
	Domain string `json:"domain"`

	// State is the state of the domain.
	State DomainState `json:"state"`

	// CreatedAt defines when this domain was created.
	CreatedAt Timestamp `json:"created_at"`
}

// Gateway is used to define a gateway used in Ignite.
type Gateway struct {
	// ID is used to define the ID of a gateway.
	ID string `json:"id"`

	// Type is the type of the gateway.
	Type GatewayType `json:"type"`

	// Name is the name of the gateway.
	Name string `json:"name"`

	// Protocol is the protocol of the gateway. This is only used on external gateways. This will be blank on internal gateways.
	Protocol GatewayProtocol `json:"protocol"`

	// DeploymentID is the ID of the deployment that this gateway is attached to.
	DeploymentID string `json:"deployment_id"`

	// CreatedAt defines when this gateway was created.
	CreatedAt Timestamp `json:"created_at"`

	// HopshDomain is the hop.sh domain that this gateway is automatically assigned. This will be blank if none is assigned.
	HopshDomain string `json:"hopsh_domain"`

	// InternalDomain is the internal domain that this gateway is automatically assigned. This will be blank if none is assigned.
	InternalDomain string `json:"internal_domain"`

	// TargetPort is the port that this gateway is targeting. This will be nil if none is assigned.
	TargetPort *int `json:"target_port"`

	// Domains is the list of domains that this gateway is assigned to.
	Domains []*Domain `json:"domains"`
}
