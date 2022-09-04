package hopgo

// Users is used to define the API methods which are part of the Users API set. Please use the Users field on
// the client made by NewClient to get an instance of this.
type Users struct {
	c *Client
}

func newUsers(c *Client) *Users {
	return &Users{c}
}
