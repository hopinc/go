package hopgo

// Pipe is used to define the API methods which are part of the Pipe API set. Please use the Pipe field on the client
// made by NewClient to get an instance of this.
type Pipe struct {
	c *Client
}
