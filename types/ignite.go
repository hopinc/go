package types

import "encoding/json"

// VolumeFormat is used to define the format of a volume.
type VolumeFormat string

const (
	// VolumeFormatExt4 defines a volume format of ext4.
	VolumeFormatExt4 VolumeFormat = "ext4"

	// VolumeFormatXFS defines a volume format of xfs.
	VolumeFormatXFS VolumeFormat = "xfs"
)

// VolumeDefinition is used to define a volume definition.
type VolumeDefinition struct {
	// FS is the filesystem of the volume.
	FS VolumeFormat `json:"fs"`

	// Size is the size of the volume.
	Size Size `json:"size"`

	// MountPath is the mount path of the volume.
	MountPath string `json:"mountpath"`
}

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

	// HopshDomainEnabled determines if the hop.sh domain is currently active.
	HopshDomainEnabled bool `json:"hopsh_domain_enabled"`

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
	// RuntimeTypeEphemeral are sort of fire and forget. Containers won't restart if they exit,
	// but they can still be terminated programmatically.
	RuntimeTypeEphemeral RuntimeType = "ephemeral"

	// RuntimeTypePersistent will restart if they exit. They can also be started and stopped programmatically.
	RuntimeTypePersistent RuntimeType = "persistent"

	// RuntimeTypeStateful is for deployments/containers can only run one container at a time, and will have a persistent volume attached.
	RuntimeTypeStateful RuntimeType = "stateful"
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
	VCPU float64 `json:"vcpu"`
	RAM  Size    `json:"ram"`
	VGPU []VGPU  `json:"vgpu"`
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

// RestartPolicy is used to define the restart policy of a deployment.
type RestartPolicy string

const (
	// RestartPolicyNever is used to define a deployment that never restarts.
	RestartPolicyNever RestartPolicy = "never"

	// RestartPolicyAlways is used to define a deployment that always restarts.
	RestartPolicyAlways RestartPolicy = "always"

	// RestartPolicyOnFailure is used to define a deployment that restarts on failure.
	RestartPolicyOnFailure RestartPolicy = "on-failure"
)

// DeploymentConfigPartial is the partial configuration of a deployment.
type DeploymentConfigPartial struct {
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

	// RestartPolicy is the restart policy for this deployment.
	RestartPolicy RestartPolicy `json:"restart_policy"`

	// Entrypoint is used to define the entrypoint for the application. Can be nil.
	Entrypoint []string `json:"entrypoint"`
}

func (x DeploymentConfigPartial) makeMap() map[string]any {
	if x.Version == "" {
		x.Version = "2022-05-17"
	}
	if x.Env == nil {
		x.Env = map[string]string{}
	}
	return map[string]any{
		"container_strategy": x.ContainerStrategy,
		"type":               x.Type,
		"version":            x.Version,
		"image":              x.Image,
		"env":                x.Env,
		"resources":          x.Resources,
		"restart_policy":     x.RestartPolicy,
		"entrypoint":         x.Entrypoint,
	}
}

// MarshalJSON is used to marshal the deployment config into JSON.
func (x DeploymentConfigPartial) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.makeMap())
}

// DeploymentConfig is used to define the configuration for a deployment.
type DeploymentConfig struct {
	// DeploymentConfigPartial is the partial configuration of a deployment that this is based on.
	DeploymentConfigPartial `json:",inline"`

	// Name is the name of the deployment.
	Name string `json:"name"`

	// Volume is the volume that this deployment is using. This can only be used when Type is RuntimeTypeStateful.
	Volume *VolumeDefinition `json:"volume,omitempty"`
}

// MarshalJSON is used to marshal the deployment config into JSON.
func (x DeploymentConfig) MarshalJSON() ([]byte, error) {
	m := x.makeMap()
	m["name"] = x.Name
	if x.Volume != nil {
		m["volume"] = x.Volume
	}
	return json.Marshal(m)
}

// DeploymentRollout is used to define the rollout of a deployment.
type DeploymentRollout struct {
	// Count is the number of containers which are being recreated.
	Count int `json:"count"`

	// CreatedAt is the time that the rollout was created at.
	CreatedAt Timestamp `json:"created_at"`

	// DeploymentID is the ID of the deployment that this rollout is for.
	DeploymentID string `json:"deployment_id"`

	// ID is the ID of the rollout.
	ID string `json:"id"`

	// State is the state of the rollout.
	State RolloutState `json:"status"`

	// Build is the build that this rollout is for. This can be nil.
	Build *Build `json:"build"`
}

// BuildMethod is the method used to build a container.
type BuildMethod string

const (
	// BuildMethodGitHub is used to define a build method that was triggered by GitHub.
	BuildMethodGitHub BuildMethod = "github"

	// BuildMethodCLI is used to define a build method that was triggered by the CLI.
	BuildMethodCLI BuildMethod = "cli"
)

// BuildAuthor is used to define the author of a build.
type BuildAuthor struct {
	// AvatarURL is the URL of the avatar of the author. This can be blank.
	AvatarURL string `json:"avatar_url"`

	// Username is the username of the author.
	Username string `json:"username"`
}

// BuildMetadata is the metadata for a build.
type BuildMetadata struct {
	// AccountType is the type of account that triggered the build. This can be blank.
	AccountType string `json:"account_type"`

	// Author is the author of the build. This can be nil.
	Author *BuildAuthor `json:"author"`

	// RepoID is the ID of the repository that the build was triggered from.
	RepoID int `json:"repo_id"`

	// RepoName is the name of the repository that the build was triggered from.
	RepoName string `json:"repo_name"`

	// Branch is the branch that the build was triggered from.
	Branch string `json:"branch"`

	// CommitSHA is the SHA of the commit that the build was triggered from.
	CommitSHA string `json:"commit_sha"`

	// CommitMessage is the message of the commit that the build was triggered from.
	CommitMessage string `json:"commit_msg"`

	// CommitURL is the URL of the commit that the build was triggered from. This can be blank.
	CommitURL string `json:"commit_url"`
}

// BuildState is the state of a build.
type BuildState string

const (
	// BuildStatePending is used to define a build that is pending.
	BuildStatePending BuildState = "pending"

	// BuildStateFailed is used to define a build that has failed.
	BuildStateFailed BuildState = "failed"

	// BuildStateSucceeded is used to define a build that has succeeded.
	BuildStateSucceeded BuildState = "succeeded"

	// BuildStateCancelled is used to define a build that has been canceled.
	BuildStateCancelled BuildState = "cancelled"
)

// Build is used to define the active build of a deployment.
type Build struct {
	// ID is the ID of the build.
	ID string `json:"id"`

	// DeploymentID is the ID of the deployment that this build is for.
	DeploymentID string `json:"deployment_id"`

	// Metadata is the metadata for the build. This can be nil.
	Metadata *BuildMetadata `json:"metadata"`

	// Method is the method used to build the container.
	Method BuildMethod `json:"method"`

	// StartedAt is the time that the build was started at. Is nil for no value.
	StartedAt *Timestamp `json:"started_at"`

	// FinishedAt is the time that the build finished at. Is nil for no value.
	FinishedAt *Timestamp `json:"finished_at"`

	// State is the state of the build.
	State BuildState `json:"state"`

	// Digest is the digest for the image. Can be blank for no value.
	Digest string `json:"digest"`
}

// DeploymentMetadata is the deployments metadata.
type DeploymentMetadata struct {
	// ContainerPortMappings is used to map the containers to ports.
	ContainerPortMappings map[string][]string `json:"container_port_mappings"`
}

// Deployment is used to define a deployment in Ignite.
type Deployment struct {
	// ID is the ID of the deployment.
	ID string `json:"id"`

	// Name is the name of the deployment.
	Name string `json:"name"`

	// TargetContainerCount is the number of expected containers.
	TargetContainerCount int `json:"target_container_count"`

	// ContainerCount is the number of containers that are currently running.
	ContainerCount int `json:"container_count"`

	// CreatedAt defines when this deployment was created.
	CreatedAt Timestamp `json:"created_at"`

	// Config is the configuration for this deployment.
	Config DeploymentConfigPartial `json:"config"`

	// ActiveRollout is the rollout for this deployment. Can be nil if not defined.
	ActiveRollout *DeploymentRollout `json:"active_rollout"`

	// ActiveBuild is the build for this deployment. Can be nil if not defined.
	ActiveBuild *Build `json:"active_build"`

	// Metadata is used to define any deployment metadata. Can be nil if not defined.
	Metadata *DeploymentMetadata `json:"metadata"`

	// RunningContainerCount is the amount of containers in the running state.
	RunningContainerCount int `json:"running_container_count"`
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

// ContainerMetadata is used to define the metadata for a container.
type ContainerMetadata struct {
	// LastExitCode is used to define the last exit code of the container. It is nil if the container has never exited.
	LastExitCode *int `json:"last_exit_code"`
}

// ContainerMetrics is used to define the metrics of a container.
type ContainerMetrics struct {
	// CPUUsagePercent is used to define the % usage of the CPU.
	CPUUsagePercent int `json:"cpu_usage_percent"`

	// MemoryUsagePercent is used to define the % usage of the RAM.
	MemoryUsagePercent int `json:"memory_usage_percent"`
}

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

	// Metadata is the metadata for this container.
	Metadata ContainerMetadata `json:"metadata"`

	// Volume is the volume definition for this container. This can be nil.
	Volume *VolumeDefinition `json:"volume"`

	// Metrics is used to define the container metrics. This can be nil.
	Metrics *ContainerMetrics `json:"metrics"`
}

// GatewayCreationOptions is used to define the options for creating a gateway.
type GatewayCreationOptions struct {
	// DeploymentID is the ID of the deployment that this gateway is for.
	DeploymentID string `json:"-"`

	// Name is the name of the gateway.
	Name string `json:"name"`

	// Type is the type of gateway to create, either internal or external.
	Type GatewayType `json:"type"`

	// Protocol is the protocol to use for the gateway.
	Protocol GatewayProtocol `json:"protocol"`

	// TargetPort is the port to listen on.
	TargetPort int `json:"target_port"`

	// ProjectID is the ID of the project that this gateway is for. Can be blank if using a project token.
	//
	// Deprecated: Set the project ID with client options instead.
	ProjectID string `json:"-"`
}

// LoggingLevel is used to define the logging level.
type LoggingLevel string

const (
	// LoggingLevelInfo is used to define the level of logging as informative. Stdout becomes info.
	LoggingLevelInfo LoggingLevel = "info"

	// LoggingLevelError is used to define the level of logging as an error.
	LoggingLevelError LoggingLevel = "error"
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

// IgniteDeploymentUpdateOpts is used to define the options for updating a deployment.
type IgniteDeploymentUpdateOpts struct {
	// Name is the name of the deployment. If this is not blank, it will be updated.
	Name string `json:"name,omitempty"`

	// Image is the image to use for the deployment. If this is not nil, it will be updated.
	Image *Image `json:"image,omitempty"`

	// RestartPolicy is the restart policy for the deployment. If this is not blank, it will be updated.
	RestartPolicy RestartPolicy `json:"restart_policy,omitempty"`

	// ContainerStrategy is the container strategy for the deployment. If this is not blank, it will be updated.
	ContainerStrategy ContainerStrategy `json:"container_strategy,omitempty"`

	// Type is the runtime type for the deployment. If this is not blank, it will be updated.
	Type RuntimeType `json:"type,omitempty"`

	// Resources is the resources for the deployment. If this is not nil, it will be updated.
	Resources *Resources `json:"resources,omitempty"`
}

// IgniteDeploymentPatchOpts is the old name for IgniteDeploymentUpdateOpts.
//
// Deprecated: Use IgniteDeploymentUpdateOpts instead.
type IgniteDeploymentPatchOpts = IgniteDeploymentUpdateOpts

// RolloutState is used to define the state of a rollout.
type RolloutState string

const (
	// RolloutStatePending is used to define a rollout that is pending.
	RolloutStatePending RolloutState = "pending"

	// RolloutStateFinished is used to define a rollout that has finished.
	RolloutStateFinished RolloutState = "finished"

	// RolloutStateFailed is used to define a rollout that has failed.
	RolloutStateFailed RolloutState = "failed"
)

// IgniteGatewayUpdateOpts is used to define the options for updating a gateway.
type IgniteGatewayUpdateOpts struct {
	// Name is the name of the gateway. If this is not blank, it will be updated.
	Name string `json:"name,omitempty"`

	// TargetPort is the port to listen on. If this is not 0, it will be updated.
	TargetPort int `json:"target_port,omitempty"`

	// Protocol is the protocol to use for the gateway. If this is not blank, it will be updated.
	Protocol GatewayProtocol `json:"protocol,omitempty"`
}

// HealthCheckProtocol is the type for a health check.
type HealthCheckProtocol string

const (
	// HealthCheckProtocolHTTP is used to define a health check running on HTTP.
	HealthCheckProtocolHTTP HealthCheckProtocol = "http"
)

// HealthCheckUpdateOpts is used to define update options for health check creation. All values except DeploymentID and
// HealthCheckID are optional.
type HealthCheckUpdateOpts struct {
	// DeploymentID is used to define the ID of the deployment. This must be set.
	DeploymentID string `json:"-"`

	// HealthCheckID is used to define the ID of the health check. This must be set.
	HealthCheckID string `json:"-"`

	// Protocol is the protocol that this health check will work on.
	Protocol HealthCheckProtocol `json:"protocol,omitempty"`

	// Path is the path which should be hit for the health check.
	Path string `json:"path,omitempty"`

	// Port is the port that should be hit for the health check.
	Port int `json:"port,omitempty"`

	// InitialDelay is the initial delay in health checking.
	InitialDelay Seconds `json:"initial_delay,omitempty"`

	// Interval is the interval between health checks in seconds.
	Interval Seconds `json:"interval,omitempty"`

	// Timeout is used to define the timeout in milliseconds.
	Timeout Milliseconds `json:"timeout,omitempty"`

	// MaxRetries is the maximum number of allowed retries before it is declared unhealthy.
	MaxRetries int `json:"max_retries,omitempty"`
}

// HealthCheckCreateOpts is used to define options during health check creation.
type HealthCheckCreateOpts struct {
	// DeploymentID is used to define the ID of the deployment. This must be set.
	DeploymentID string `json:"deployment_id,omitempty"`

	// Protocol is the protocol that this health check will work on. If blank, will default to "http".
	Protocol HealthCheckProtocol `json:"protocol"`

	// Path is the path which should be hit for the health check. If blank, will default to "/".
	Path string `json:"path"`

	// Port is the port that should be hit for the health check. If blank, will default to 8080.
	Port int `json:"port"`

	// InitialDelay is the initial delay in health checking. If blank, will default to 5 seconds.
	InitialDelay Seconds `json:"initial_delay"`

	// Interval is the interval between health checks in seconds. If blank, will default to 1 minute.
	Interval Seconds `json:"interval"`

	// Timeout is used to define the timeout in milliseconds. If blank, will default to 50ms.
	Timeout Milliseconds `json:"timeout"`

	// MaxRetries is the maximum number of allowed retries before it is declared unhealthy. If blank, will default to 3.
	MaxRetries int `json:"max_retries"`
}

// HealthCheckType defines the type of the health check.
type HealthCheckType string

const (
	// HealthCheckTypeReadiness defines a readiness type.
	HealthCheckTypeReadiness HealthCheckType = "readiness"

	// HealthCheckTypeLiveness defines a liveness type.
	HealthCheckTypeLiveness HealthCheckType = "liveness"
)

// HealthCheck is used to define the created health check.
type HealthCheck struct {
	// Inlines the options since they are also used here.
	HealthCheckCreateOpts `json:",inline"`

	// ID defines the ID of the health check.
	ID string `json:"id"`

	// CreatedAt defines when the health check was created.
	CreatedAt Timestamp `json:"created_at"`

	// Type defines the type of the health check.
	Type HealthCheckType `json:"type"`
}

// HealthCheckStatus is the type of the health check status.
type HealthCheckStatus string

const (
	// HealthCheckStatusSucceeded is used to define a health check that succeeded.
	HealthCheckStatusSucceeded HealthCheckStatus = "succeeded"

	// HealthCheckStatusFailed is used to define a health check that failed.
	HealthCheckStatusFailed HealthCheckStatus = "failed"

	// HealthCheckStatusPending is used to define a health check that is pending.
	HealthCheckStatusPending HealthCheckStatus = "pending"
)

// HealthCheckState is used to define the state of a health check.
type HealthCheckState struct {
	// DeploymentID defines the deployment ID this relates to.
	DeploymentID string `json:"deployment_id"`

	// ContainerID defines the container ID this relates to.
	ContainerID string `json:"container_id"`

	// HealthCheckID is used to define the ID of the health check this relates to.
	HealthCheckID string `json:"health_check_id"`

	// State is used to define the health check status.
	State HealthCheckStatus `json:"state"`

	// NextCheck defines the timestamp of the next check.
	NextCheck Timestamp `json:"next_check"`

	// CreatedAt defines when this health check was created.
	CreatedAt Timestamp `json:"created_at"`
}

// DeploymentStorageSize is used to define the information about the build cache.
type DeploymentStorageSize struct {
	// ProvisionedSize is the amount of storage in MB that is provisioned.
	ProvisionedSize int `json:"provisioned_size"`

	// UsedSize is the amount of storage in MB that is used for the build cache.
	UsedSize int `json:"used_size"`
}

// DeploymentStorageInfo is used to define deployment information about storage.
type DeploymentStorageInfo struct {
	// Volume is used to define the storage information. Can be nil.
	Volume *DeploymentStorageSize `json:"volume"`

	// BuildCache is used to define the build cache storage information. Can be nil.
	BuildCache *DeploymentStorageSize `json:"build_cache"`
}
