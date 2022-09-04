package hopgo

// Ignite is used to define the API methods which are part of the Ignite API set. Please use the Ignite field on
// the client made by NewClient to get an instance of this.
type Ignite struct {
	c *Client
}

func newIgnite(c *Client) *Ignite {
	return &Ignite{c}
}
