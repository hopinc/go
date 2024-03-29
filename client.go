package hop

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"go.hop.io/sdk/types"
	"moul.io/http2curl"
)

var userAgent = "hop-go/" + Version + " (go/" + runtime.Version() + ")"

// DefaultAPIBase is used to define the default base API URL.
const DefaultAPIBase = "https://api.hop.io/v1"

// IDPrefixes are allowed ID prefixes.
var IDPrefixes = []string{
	"user", "project", "pm", "role", "pi", "ptk", "pat", "container", "pipe_room", "deployment", "bearer",
	"ptkid", "secret", "gateway", "domain", "leap_token", "build", "rollout", "health_check", "session",
}

// ValidateToken is used to validate a authentication token. Returns an error if the token is invalid.
func ValidateToken(authorization string) (string, error) {
	for _, v := range IDPrefixes {
		if strings.HasPrefix(authorization, v+"_") {
			return v, nil
		}
	}
	return "", types.InvalidToken("invalid authorization token prefix: " + authorization)
}

// Client is used to define the Hop client. Please use the NewClient function to create this!
type Client struct {
	httpClient    *http.Client
	authorization string
	tokenType     string
	isTest        bool
	opts          []ClientOption

	Pipe     *ClientCategoryPipe
	Projects *ClientCategoryProjects
	Ignite   *ClientCategoryIgnite
	Users    *ClientCategoryUsers
	Registry *ClientCategoryRegistry
	Channels *ClientCategoryChannels
}

// NewClient is used to make a new Hop client.
func NewClient(authorization string, opts ...ClientOption) (*Client, error) {
	prefix, err := ValidateToken(authorization)
	if err != nil {
		return nil, err
	}
	if prefix != "bearer" && prefix != "pat" && prefix != "ptk" {
		return nil, types.InvalidToken("invalid authorization token prefix: " + prefix)
	}

	var c Client
	c = Client{
		httpClient:    &http.Client{},
		opts:          opts,
		authorization: authorization,
		tokenType:     prefix,
		Pipe:          newPipe(&c),
		Projects:      newProjects(&c),
		Ignite:        newIgnite(&c),
		Users:         newUsers(&c),
		Registry:      newRegistry(&c),
		Channels:      newChannels(&c),
	}
	return &c, nil
}

// AddClientOptions is used to add any client options. These will be added before any others specified in functions.
func (c *Client) AddClientOptions(opts ...ClientOption) {
	c.opts = append(c.opts, opts...)
}

// SetAPIBase is used to set the base API URL. This is probably something you do not need to use, however it is useful
// in testing the SDK. The base URL contains the domain and ends with /v1.
//
// Deprecated: Use AddClientOptions with WithCustomAPIURL instead.
func (c *Client) SetAPIBase(apiBase string) *Client {
	c.AddClientOptions(WithCustomAPIURL(apiBase))
	return c
}

// PlainText is a special case for a body that should be sent as text/plain.
type PlainText []byte

// ClientArgs is used to define the arguments from the function to the client.
type ClientArgs struct {
	// Method is used to define the request method.
	Method string

	// Path is used to define the path (with the URL suffix including /api/vX removed).
	Path string

	// ResultKey is the key inside the "data" part of the JSON where the data is.
	ResultKey string

	// Query is any HTTP query parameters that should be sent.
	Query map[string]string

	// Body is the body to send. This should be json marshalled except for the PlainText type.
	Body any

	// Result is a pointer to where the result should be unmarshalled. Note that if this is nil, the
	// result body should be ignored.
	// *SPECIAL CASE:* If PassRequest is not nil, ignore this field and pass it to that.
	Result any

	// Ignore404 is whether a 404 should not be treated as an error.
	Ignore404 bool

	// PassRequest is a function used to pass a OK request to instead of closing it.
	// This is used for the filesystem since we need to pass around a context that constantly pulls more
	// information in some cases.
	PassRequest func(r *http.Response)
}

type responseBody struct {
	Data json.RawMessage `json:"data"`
}

func (c *Client) getTokenType() string { return c.tokenType }

// Transforms the client options into the processed object.
func (c *Client) processOpts(opts []ClientOption) ProcessedClientOpts {
	o := ProcessedClientOpts{CurlDebugWriters: []io.Writer{}}
	process := func(v any) {
		switch x := v.(type) {
		case projectIdOption:
			o.ProjectID = x.projectId
		case apiUrlOption:
			o.CustomAPIURL = x.apiBase
		case curlWriterOption:
			o.CurlDebugWriters = append(o.CurlDebugWriters, x.w)
		case customHandlerOption:
			o.CustomHandler = x.fn
		}
	}
	for _, v := range c.opts {
		process(v)
	}
	for _, v := range opts {
		process(v)
	}
	return o
}

// Resolves the project ID from the client options. Will be a blank string if this is not specified.
func (c *Client) getProjectId(opts []ClientOption) string {
	projectId := ""
	for _, v := range c.opts {
		if x, ok := v.(projectIdOption); ok {
			projectId = x.projectId
		}
	}
	for _, v := range opts {
		if x, ok := v.(projectIdOption); ok {
			projectId = x.projectId
		}
	}
	return projectId
}

func (c *Client) setRequestHeaders(req *http.Request, r io.Reader, textPlain bool) {
	req.Header.Set("Authorization", c.authorization)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if r != nil {
		// This means we have a body of some description. What content type should we use?
		if textPlain {
			req.Header.Set("Content-Type", "text/plain")
		} else {
			req.Header.Set("Content-Type", "application/json")
		}
	}
}

// Does the specified HTTP request.
func (c *Client) do(ctx context.Context, a ClientArgs, clientOpts []ClientOption) error { //nolint:funlen,gocognit,gocyclo,cyclop
	// Transform the client options.
	processedOpts := c.processOpts(clientOpts)

	// Handle client overrides.
	if processedOpts.CustomHandler != nil {
		return processedOpts.CustomHandler(ctx, a, processedOpts)
	}

	// Handle getting the body bytes.
	var body []byte
	textPlain := false
	if a.Method != "GET" && a.Body != nil {
		switch x := a.Body.(type) {
		case PlainText:
			textPlain = true
			body = x
		case []byte:
			body = x
		default:
			bodyBytes, err := json.Marshal(a.Body)
			if err != nil {
				return err
			}
			body = bodyBytes
		}
	}

	// Add project ID to the query if it is specified.
	if processedOpts.ProjectID != "" {
		if a.Query == nil {
			a.Query = map[string]string{}
		}
		a.Query["project"] = processedOpts.ProjectID
	}

	// Create the request.
	suffix := ""
	if a.Query != nil {
		suffix = "?"
		first := true

		addChunk := func(k, v string) {
			chunk := ""
			if first {
				first = false
			} else {
				chunk = "&"
			}
			suffix += chunk + url.QueryEscape(k) + "=" + url.QueryEscape(v)
		}

		if c.isTest {
			// Order all the keys for unit testing reasons.
			keys := make([]string, len(a.Query))
			i := 0
			for k := range a.Query {
				keys[i] = k
				i++
			}
			sort.Strings(keys)
			for _, k := range keys {
				addChunk(k, a.Query[k])
			}
		} else {
			// Just proceed as usual.
			for k, v := range a.Query {
				addChunk(k, v)
			}
		}
	}
	apiBase := processedOpts.CustomAPIURL
	if apiBase == "" {
		// Revert to the default.
		apiBase = DefaultAPIBase
	}

	if len(processedOpts.CurlDebugWriters) != 0 {
		// If curl debugging is on, we make the request structure twice. This is because NewRequestWithContext
		// takes a reader, and we do not want to pollute that.
		var r io.Reader
		if body != nil {
			r = bytes.NewReader(body)
		}
		curlReq, err := http.NewRequest(a.Method, apiBase+a.Path+suffix, r) //nolint:noctx // Built for curl handler.
		if err != nil {
			return err
		}
		c.setRequestHeaders(curlReq, r, textPlain)

		// Convert the request to a curl command.
		var curl *http2curl.CurlCommand
		curl, err = http2curl.GetCurlCommand(curlReq)
		if err != nil {
			return err
		}
		curlB := []byte(curl.String() + "\n")

		// Send the curl request to all the writers.
		for _, v := range processedOpts.CurlDebugWriters {
			if _, err = v.Write(curlB); err != nil {
				return err
			}
		}
	}

	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, a.Method, apiBase+a.Path+suffix, r)
	if err != nil {
		return err
	}
	c.setRequestHeaders(req, r, textPlain)

	// Do the request.
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	// The closure *can* be cancelled later in the chain, lets watch for this.
	closeBody := true
	defer func() {
		if closeBody {
			_ = res.Body.Close()
		}
	}()

	// If this is a 4xs or 5xx, handle the error.
	if res.StatusCode >= 400 && 599 >= res.StatusCode {
		// If this is a 404, check if it is a special case before jumping to the error handler.
		if res.StatusCode != 404 || !a.Ignore404 {
			return handleErrors(res)
		}
	}

	// Handle if we should pass off the request.
	if a.PassRequest != nil {
		a.PassRequest(res)
		closeBody = false
		return nil
	}

	// Handle if we should process the data.
	if a.Result != nil {
		var b []byte
		if b, err = io.ReadAll(res.Body); err != nil {
			// Failed to read the body.
			return err
		}
		var x responseBody
		if err = json.Unmarshal(b, &x); err != nil {
			// Unable to unmarshal the body.
			return err
		}
		b = x.Data

		if a.ResultKey != "" {
			// Get the json.RawMessage for the specific key.
			var m map[string]json.RawMessage
			if err = json.Unmarshal(b, &m); err != nil {
				return err
			}
			var ok bool
			if b, ok = m[a.ResultKey]; !ok {
				// The key specified was not actually valid.
				return errors.New("api response error: key was not in response - please report this to " +
					"the go-hop github repository")
			}
		}

		if err = json.Unmarshal(b, a.Result); err != nil {
			return err
		}
	}

	// Success! No errors!
	return nil
}

type clientDoer interface {
	do(ctx context.Context, a ClientArgs, opts []ClientOption) error
	getTokenType() string
	getProjectId([]ClientOption) string
}

// Paginator is used to create a way to access paginated API routes.
type Paginator[T any] struct {
	c         clientDoer
	pageIndex int
	count     int
	total     int // should be set to -1 on init.
	mu        sync.Mutex

	// Defines things required for the offset strategy.
	offsetStrat bool
	limit       int

	path      string
	resultKey string
	sortBy    string
	orderBy   string
	query     map[string]string
}

func unJsonInt(j json.RawMessage) (int, error) {
	var x int
	err := json.Unmarshal(j, &x)
	return x, err
}

// Next is used to get the next page.
func (p *Paginator[T]) Next(ctx context.Context, opts ...ClientOption) ([]T, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.total != -1 && p.count >= p.total {
		return nil, types.StopIteration
	}

	query := map[string]string{}
	for k, v := range p.query {
		query[k] = v
	}
	if p.offsetStrat {
		query["offset"] = strconv.Itoa(p.count)
		query["limit"] = strconv.Itoa(p.limit)
	} else {
		query["page"] = strconv.Itoa(p.pageIndex + 1)
	}
	if p.orderBy == "" {
		query["orderBy"] = "asc"
	} else {
		query["orderBy"] = p.orderBy
	}
	if p.sortBy != "" {
		query["sortBy"] = p.sortBy
	}

	var m map[string]json.RawMessage
	if err := p.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      p.path,
		Query:     query,
		Result:    &m,
		Ignore404: false,
	}, opts); err != nil {
		return nil, err
	}

	if totalCount, ok := m["total_count"]; ok {
		// We have a count to go by. This probably means we are not using the offset strategy.
		var err error
		p.total, err = unJsonInt(totalCount)
		if err != nil {
			return nil, err
		}
	}

	var a []T
	if err := json.Unmarshal(m[p.resultKey], &a); err != nil {
		return nil, err
	}
	if len(a) == 0 {
		// Stop pagination here.
		return nil, types.StopIteration
	}
	p.count += len(a)
	if !p.offsetStrat {
		// Add 1 to pages since this is not using the offset strategy.
		p.pageIndex++
	}
	return a, nil
}

// ForChunk is basically the shorthand for calling a function everytime there is a new result. Any errors are passed to
// the root error result.
func (p *Paginator[T]) ForChunk(ctx context.Context, f func([]T) error, opts ...ClientOption) error {
	for a, err := p.Next(ctx, opts...); err != types.StopIteration; a, err = p.Next(ctx) {
		if err != nil {
			return err
		}
		if err = f(a); err != nil {
			return err
		}
	}
	return nil
}
