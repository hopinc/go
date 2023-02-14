package types

// User is used to define a user.
type User struct {
	// ID is the ID of the user.
	ID string `json:"id"`

	// Name is the name of the user.
	Name string `json:"name"`

	// Username is the username of the user.
	Username string `json:"username"`

	// Email is the email of the user.
	Email string `json:"email"`
}

// SelfUser is used to define the self user.
type SelfUser struct {
	User `json:",inline"`

	// EmailVerified is set to true if the user has verified their email.
	EmailVerified bool `json:"email_verified"`

	// TOTPEnabled is set to true if the user has enabled TOTP authentication.
	TOTPEnabled bool `json:"totp_enabled"`

	// WebauthnEnabled is set to true if the user hss enabled webauthn authentication.
	WebauthnEnabled bool `json:"webauthn_enabled"`

	// MFAEnabled is set to true if the user has enabled MFA.
	MFAEnabled bool `json:"mfa_enabled"`

	// Admin defines if the user is an admin.
	Admin bool `json:"admin"`
}

// UserMeInfo is the payload returned fcr all information about the current user.
type UserMeInfo struct {
	// Projects is the list of projects the user is a member of.
	Projects []*Project `json:"projects"`

	// User is the user.
	User SelfUser `json:"user"`

	// ProjectMemberRoleMap is a map of project ID to project member role.
	ProjectMemberRoleMap map[string]*ProjectRole `json:"project_member_role_map"`

	// LeapToken is the users Leap token. Can be blank.
	LeapToken string `json:"leap_token"`
}

// UserPat is used to define a personal access token.
type UserPat struct {
	// ID is the ID of the personal access token.
	ID string `json:"id"`

	// Name is the name of the personal access token.
	Name string `json:"name"`

	// PAT is a partially censored personal access token (unless this is from a creation).
	PAT string `json:"pat"`

	// CreatedAt is when the personal access token was created.
	CreatedAt Timestamp `json:"created_at"`
}
