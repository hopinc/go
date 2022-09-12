package hop

// ClientOption is used to define am option that the client will consume when it is ran.
// You should not and cannot implement this interface. This is only designed as a signature for other options.
type ClientOption interface {
	youCannotImplementClientOption()
}

type baseClientOption struct{}

func (baseClientOption) youCannotImplementClientOption() {}

type projectIdOption struct {
	baseClientOption

	projectId string
}

// WithProjectID Is used to return a ClientOption for calling the API with a specific project ID.
func WithProjectID(projectId string) ClientOption {
	return projectIdOption{projectId: projectId}
}
