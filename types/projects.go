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
	VCPU int `json:"vcpu"`
	RAM  int `json:"ram"`
}

// QuotaUsage is used to define the usage for a specified quota.
type QuotaUsage struct {
	VCPU int `json:"vcpu"`
	RAM  int `json:"ram"`
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
