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

	"github.com/hopinc/hop-go/types"
)

var userAgent = "hop-go/" + Version + " (go/" + runtime.Version() + ")"

// DefaultAPIBase is used to define the default base API URL.
const DefaultAPIBase = "https://api.hop.io/v1"

// IDPrefixes are allowed ID prefixes.
var IDPrefixes = []string{
	"user", "project", "pm", "role", "pi", "ptk", "pat", "container", "pipe_room", "deployment", "bearer",
	"ptkid", "secret", "gateway", "domain", "leap_token", "build",
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
	apiBase       string
	isTest        bool
	opts          []ClientOption

	Pipe     *ClientCategoryPipe
	Projects *ClientCategoryProjects
	Ignite   *ClientCategoryIgnite
	Users    *ClientCategoryUsers
	Registry *ClientCategoryRegistry
	Channels *ClientCategoryChannels
}

// Runs through all options in the client and then any additional options passed through.
func (c *Client) forOption(f func(any), opts []ClientOption) {
	if c.opts != nil {
		for _, v := range c.opts {
			f(v)
		}
	}
	for _, v := range opts {
		f(v)
	}
}

// NewClient is used to make a new Hop client.
func NewClient(authorization string) (*Client, error) {
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
func (c *Client) SetAPIBase(apiBase string) *Client {
	if !strings.HasPrefix(apiBase, "http://") && !strings.HasPrefix(apiBase, "https://") {
		apiBase = "https://" + apiBase
	}
	if strings.HasSuffix(apiBase, "/") {
		// Remove the trailing slash.
		apiBase = apiBase[:len(apiBase)-1]
	}
	c.apiBase = apiBase
	return c
}

type plainText []byte

type clientArgs struct {
	method    string
	path      string
	resultKey string
	query     map[string]string
	body      any
	result    any
	ignore404 bool
}

type responseBody struct {
	Data json.RawMessage `json:"data"`
}

func (c *Client) getTokenType() string { return c.tokenType }

// Resolves the project ID from the client options. Will be a blank string if this is not specified.
func (c *Client) getProjectId(opts []ClientOption) string {
	projectId := ""
	c.forOption(func(x any) {
		if x, ok := x.(projectIdOption); ok {
			projectId = x.projectId
		}
	}, opts)
	return projectId
}

// Does the specified HTTP request.
func (c *Client) do(ctx context.Context, a clientArgs, clientOpts []ClientOption) error {
	// Handle getting the body bytes.
	var r io.Reader
	textPlain := false
	if a.method != "GET" && a.body != nil {
		switch x := a.body.(type) {
		case plainText:
			textPlain = true
			r = bytes.NewReader(x)
		case []byte:
			r = bytes.NewReader(x)
		default:
			bodyBytes, err := json.Marshal(a.body)
			if err != nil {
				return err
			}
			r = bytes.NewReader(bodyBytes)
		}
	}

	// Add project ID to the query if it is specified.
	if projectId := c.getProjectId(clientOpts); projectId != "" {
		if a.query == nil {
			a.query = map[string]string{}
		}
		a.query["project"] = projectId
	}

	// Create the request.
	suffix := ""
	if a.query != nil {
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
			keys := make([]string, len(a.query))
			i := 0
			for k := range a.query {
				keys[i] = k
				i++
			}
			sort.Strings(keys)
			for _, k := range keys {
				addChunk(k, a.query[k])
			}
		} else {
			// Just proceed as usual.
			for k, v := range a.query {
				addChunk(k, v)
			}
		}
	}
	apiBase := c.apiBase
	if apiBase == "" {
		apiBase = DefaultAPIBase
	}
	req, err := http.NewRequestWithContext(ctx, a.method, apiBase+a.path+suffix, r)
	if err != nil {
		return err
	}
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

	// Do the request.
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// If this is a 4xs or 5xx, handle the error.
	if res.StatusCode >= 400 && 599 >= res.StatusCode {
		// If this is a 404, check if it is a special case before jumping to the error handler.
		if res.StatusCode != 404 || !a.ignore404 {
			return handleErrors(res)
		}
	}

	// Handle if we should process the data.
	if a.result != nil {
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

		if a.resultKey != "" {
			// Get the json.RawMessage for the specific key.
			var m map[string]json.RawMessage
			if err = json.Unmarshal(b, &m); err != nil {
				return err
			}
			var ok bool
			if b, ok = m[a.resultKey]; !ok {
				// The key specified was not actually valid.
				return errors.New("api response error: key was not in response - please report this to " +
					"the go-hop github repository")
			}
		}

		if err = json.Unmarshal(b, a.result); err != nil {
			return err
		}
	}

	// Success! No errors!
	return nil
}

type clientDoer interface {
	do(ctx context.Context, a clientArgs, opts []ClientOption) error
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
	if err := p.c.do(ctx, clientArgs{
		method:    "GET",
		path:      p.path,
		query:     query,
		result:    &m,
		ignore404: false,
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
