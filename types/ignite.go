package types

import "encoding/json"

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

// ContainerStrategy is the strategy used to scale a container.
type ContainerStrategy string

const (
	// ContainerStrategyManual is used to define a container that is manually scaled.
	ContainerStrategyManual ContainerStrategy = "manual"
)

// RuntimeType is used to define the type of runtime.
type RuntimeType string

const (
	// RuntimeTypeEphemeral are sort of fire and forget. Containers won't restart if they exit but they can still be terminated programmatically.
	RuntimeTypeEphemeral RuntimeType = "ephemeral"

	// RuntimeTypePersistent will restart if they exit. They can also be started and stopped programmatically.
	RuntimeTypePersistent RuntimeType = "persistent"
)

// DockerAuth is used to define the authentication information for a Docker registry.
type DockerAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ImageGHInfo is used to define the information about a GitHub image.
type ImageGHInfo struct {
	RepoID   int    `json:"repo_id"`
	FullName string `json:"full_name"`
	Branch   string `json:"branch"`
}

// Image is used to define an image in Ignite.
type Image struct {
	// Name is the name of the image. Will be blank if there's no image name.
	Name string `json:"name"`

	// Auth is the authentication information for the image. Will be nil if there's no authentication information.
	Auth *DockerAuth `json:"auth"`

	// GithubRepo is the GitHub repository that this image is from. Will be nil if there's no GitHub repository.
	GithubRepo *ImageGHInfo `json:"github_repo"`
}

// GPUType is the type of GPU.
type GPUType string

const (
	// GPUTypeA400 is used to define an A400 GPU.
	GPUTypeA400 GPUType = "a400"
)

// VGPU is used to define a virtual GPU.
type VGPU struct {
	// Type is the type of the virtual GPU.
	Type GPUType `json:"type"`

	// Count is the number of GPU's to allocate of this type.
	Count int `json:"count"`
}

// Resources is used to define the resources used by a deployment.
type Resources struct {
	VCPU int    `json:"vcpu"`
	RAM  Size   `json:"ram"`
	VGPU []VGPU `json:"vgpu"`
}

// MarshalJSON is used to marshal the resources into JSON.
func (r Resources) MarshalJSON() ([]byte, error) {
	if r.VGPU == nil {
		r.VGPU = []VGPU{}
	}
	return json.Marshal(map[string]any{
		"vcpu": r.VCPU,
		"ram":  r.RAM,
		"vgpu": r.VGPU,
	})
}

// DeploymentConfig is used to define the configuration for a deployment.
type DeploymentConfig struct {
	// Name is the name of the deployment.
	Name string `json:"name"`

	// ContainerStrategy is the strategy used to scale a container.
	ContainerStrategy ContainerStrategy `json:"container_strategy"`

	// Type is used to define the type of this deployment.
	Type RuntimeType `json:"type"`

	// Version is the version of the deployment.
	Version string `json:"version"`

	// Image is the Docker image config for this deployment.
	Image Image `json:"image"`

	// Env is the environment variables for this deployment.
	Env map[string]string `json:"env"`

	// Resources is the resources for this deployment.
	Resources Resources `json:"resources"`
}

// MarshalJSON is used to marshal the deployment config into JSON.
func (c DeploymentConfig) MarshalJSON() ([]byte, error) {
	if c.Version == "" {
		c.Version = "2022-05-17"
	}
	if c.Env == nil {
		c.Env = map[string]string{}
	}
	return json.Marshal(map[string]any{
		"name":               c.Name,
		"container_strategy": c.ContainerStrategy,
		"type":               c.Type,
		"version":            c.Version,
		"image":              c.Image,
		"env":                c.Env,
		"resources":          c.Resources,
	})
}

// Deployment is used to define a deployment in Ignite.
type Deployment struct {
	// ID is the ID of the deployment.
	ID string `json:"id"`

	// Name is the name of the deployment.
	Name string `json:"name"`

	// ContainerCount is the number of containers that are currently running.
	ContainerCount int `json:"container_count"`

	// CreatedAt defines when this deployment was created.
	CreatedAt Timestamp `json:"created_at"`

	// Config is the configuration for this deployment.
	Config *DeploymentConfig `json:"config"`
}

// Region is used to define a Hop datacenter region.
type Region string

const (
	// RegionUSEast1 is used to define the US East 1 region.
	RegionUSEast1 Region = "us-east-1"
)

// ContainerUptime is the structure that contains information about a containers uptime.
type ContainerUptime struct {
	// LastStart is the last time the container was started.
	LastStart Timestamp `json:"last_start"`
}

// ContainerState is used to define the current status of a container.
type ContainerState string

const (
	// ContainerStatePending is used to define a container that is pending.
	ContainerStatePending ContainerState = "pending"

	// ContainerStateRunning is used to define a container that is running.
	ContainerStateRunning ContainerState = "running"

	// ContainerStateStopped is used to define a container that is stopped.
	ContainerStateStopped ContainerState = "stopped"

	// ContainerStateFailed is used to define a container that has failed (e.g. exited with a non-zero exit code).
	ContainerStateFailed ContainerState = "failed"

	// ContainerStateTerminating is used to define a container that is terminating.
	ContainerStateTerminating ContainerState = "terminating"

	// ContainerStateExited is used to define a container that has exited.
	ContainerStateExited ContainerState = "exited"
)

// Container is used to define a container in Ignite.
type Container struct {
	// ID is the ID of the container.
	ID string `json:"id"`

	// CreatedAt defines when this container was created.
	CreatedAt Timestamp `json:"created_at"`

	// Region is the region that this container is running in.
	Region Region `json:"region"`

	// Uptime is used to define uptime/downtime information for this container.
	Uptime ContainerUptime `json:"uptime"`

	// Type is the runtime type of this container.
	Type RuntimeType `json:"type"`

	// InternalIP is the internal IP address of this container.
	InternalIP string `json:"internal_ip"`

	// DeploymentID is the ID of the deployment that this container is a part of.
	DeploymentID string `json:"deployment_id"`

	// State is the state of this container.
	State ContainerState `json:"state"`
}

// GatewayCreationOptions is used to define the options for creating a gateway.
type GatewayCreationOptions struct {
	// ProjectID is the ID of the project that this gateway is for. Can be blank if using a project token.
	ProjectID string `json:"-"`

	// DeploymentID is the ID of the deployment that this gateway is for.
	DeploymentID string `json:"-"`

	// Type is the type of gateway to create, either internal or external.
	Type GatewayType `json:"type"`

	// Protocol is the protocol to use for the gateway.
	Protocol GatewayProtocol `json:"protocol"`

	// TargetPort is the port to listen on.
	TargetPort int `json:"target_port"`
}

// LoggingLevel is used to define the logging level.
type LoggingLevel string

const (
	// LoggingLevelInfo is used to define the level of logging as informative. Stdout becomes info.
	LoggingLevelInfo LoggingLevel = "info"

	// LoggingLevelStderr is used to define the level of logging as an error.
	LoggingLevelStderr LoggingLevel = "stderr"
)

// ContainerLog is used to define a container log message.
type ContainerLog struct {
	// Timestamp is the timestamp of the log message.
	Timestamp Timestamp `json:"timestamp"`

	// Message is the log message.
	Message string `json:"message"`

	// Nonce is the ID of the document in Elasticsearch. This can be safely used to map state.
	Nonce string `json:"nonce"`

	// Level is the logging level.
	Level LoggingLevel `json:"level"`
}
