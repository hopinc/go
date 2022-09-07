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

// UserMeInfo is the payload returned fcr all information about the current user.
type UserMeInfo struct {
	// Projects is the list of projects the user is a member of.
	Projects []*Project `json:"projects"`

	// User is the user.
	User User `json:"user"`

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
