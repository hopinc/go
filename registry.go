package gop

// Registry is used to define the API methods which are part of the Registry API set. Please use the Registry field on
// the client made by NewClient to get an instance of this.
type Registry struct {
	c *Client
}
