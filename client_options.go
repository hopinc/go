package hop

import (
	"context"
	"io"
	"strings"
)

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

type apiUrlOption struct {
	baseClientOption

	apiBase string
}

// WithCustomAPIURL Is used to return a ClientOption for calling the API with a specific API URL.
// By default, this is https://api.hop.io/v1.
func WithCustomAPIURL(apiBase string) ClientOption {
	if !strings.HasPrefix(apiBase, "http://") && !strings.HasPrefix(apiBase, "https://") {
		apiBase = "https://" + apiBase
	}
	apiBase = strings.TrimSuffix(apiBase, "/")
	return apiUrlOption{apiBase: apiBase}
}

type curlWriterOption struct {
	baseClientOption

	w io.Writer
}

// WithCurlDebugWriter is used to write what the curl command for the request specified would be and a new line to an
// io.Writer.
func WithCurlDebugWriter(w io.Writer) ClientOption {
	return curlWriterOption{w: w}
}

type customHandlerOption struct {
	baseClientOption

	fn CustomHandler
}

// CustomHandler is used to define a custom handler for Hop requests.
type CustomHandler func(ctx context.Context, a ClientArgs, opts ProcessedClientOpts) error

// WithCustomHandler is used to define a custom Hop request handler.
func WithCustomHandler(fn CustomHandler) ClientOption {
	return customHandlerOption{fn: fn}
}

// ProcessedClientOpts is the result of all the client options that were passed in.
type ProcessedClientOpts struct {
	// ProjectID is the project iD this is relating to. Blank if not set.
	ProjectID string

	// CustomAPIURL is the last API URL that was swt. Blank if not set.
	CustomAPIURL string

	// CurlDebugWriters are writers that the curl command and new line of a request should be sent to.
	CurlDebugWriters []io.Writer

	// CustomHandler is used to define the custom handler for Hop requests. Nil if not set.
	CustomHandler CustomHandler
}
