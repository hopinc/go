package types

// ProjectTier is the type for a tier of a project.
type ProjectTier string

const (
	// ProjectTierPaid is used to define a project which is in the paid tier.
	ProjectTierPaid ProjectTier = "paid"

	// ProjectTierFree is used to define a project which is in the free tier.
	ProjectTierFree ProjectTier = "free"
)

// ProjectType is the type of the project.
type ProjectType string

const (
	// ProjectTypeRegular is used to define a regular project.
	ProjectTypeRegular ProjectType = "regular"

	// ProjectTypePersonal is a personal project that is created when you register an account.
	ProjectTypePersonal ProjectType = "personal"
)

// DefaultQuotas is used to define the default quotas for a hardware resource.
type DefaultQuotas struct {
	VCPU float64 `json:"vcpu"`
	RAM  float64 `json:"ram"`
}

// QuotaUsage is used to define the usage for a specified quota.
type QuotaUsage struct {
	VCPU float64 `json:"vcpu"`
	RAM  float64 `json:"ram"`
}

// Project defines the structure of a project.
type Project struct {
	// ID is used to define the ID of a project.
	ID string `json:"id"`

	// Name is the name of a project.
	Name string `json:"name"`

	// Tier is the tier of this project.
	Tier ProjectTier `json:"tier"`

	// CreatedAt defines when this project was created.
	CreatedAt Timestamp `json:"created_at"`

	// Icon is used to define the icon for this project.
	Icon *string `json:"icon"`

	// Namespace is the registry namespace for this project.
	Namespace string `json:"namespace"`

	// Type is the type of the project.
	Type ProjectType `json:"type"`

	// DefaultQuotas is used to define the default quotas for this project.
	DefaultQuotas DefaultQuotas `json:"default_quotas"`

	// QuotaOverrides is used to define any overrides to quotas.
	QuotaOverrides map[string]int `json:"quota_overrides"`

	// QuotaUsage is the quota usage for this project.
	QuotaUsage QuotaUsage `json:"quota_usage"`
}

// ProjectToken is used to define the structure of a project token.
type ProjectToken struct {
	// ID is the ID of the project token.
	ID string `json:"id"`

	// Token is part of the key value. This will likely have half of the key obfuscated.
	Token string `json:"token"`

	// CreatedAt is when this project token was created.
	CreatedAt Timestamp `json:"created_at"`
}

// ProjectPermission is used to define a permission for a project token.
type ProjectPermission string

const (
	// ProjectPermissionManageRegistry is used to define the permission to manage the registry.
	ProjectPermissionManageRegistry ProjectPermission = "MANAGE_REGISTRY"

	// ProjectPermissionManageMembers is used to define the permission to manage members.
	ProjectPermissionManageMembers ProjectPermission = "MANAGE_MEMBERS"

	// ProjectPermissionManagePipe is used to define the permission to manage the pipe.
	ProjectPermissionManagePipe ProjectPermission = "MANAGE_PIPE"

	// ProjectPermissionManageChannels is used to define the permission to manage channels.
	ProjectPermissionManageChannels ProjectPermission = "MANAGE_CHANNELS"

	// ProjectPermissionManageDeployments is used to define the permission to manage deployments.
	ProjectPermissionManageDeployments ProjectPermission = "MANAGE_DEPLOYMENTS"
)

// ProjectRole is used to define a role of a member in a project.
type ProjectRole struct {
	// ID is the ID of the role.
	ID string `json:"id"`

	// Name is the name of the role.
	Name string `json:"name"`

	// Flags is the flags for this role.
	Flags int `json:"flags"`
}

// ProjectMember is used to define a member of a project.
type ProjectMember struct {
	// ID is the ID of the member.
	ID string `json:"id"`

	// Name is the name of the member.
	Name string `json:"name"`

	// Username is the username of the member.
	Username string `json:"username"`

	// Role is the role of the member.
	Role ProjectRole `json:"role"`

	// JoinedAt is when the member joined the project.
	JoinedAt Timestamp `json:"joined_at"`
}

// ProjectSecret is used to define a secret for a project.
type ProjectSecret struct {
	// ID is the ID of the secret.
	ID string `json:"id"`

	// Name is the name of the secret.
	Name string `json:"name"`

	// Digest is the string hash of the secret.
	Digest string `json:"digest"`

	// CreatedAt is when the secret was created.
	CreatedAt Timestamp `json:"created_at"`
}
