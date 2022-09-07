package hopgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/hopinc/hop-go/types"
)

// DefaultAPIBase is used to define the default base API URL.
const DefaultAPIBase = "https://api.hop.io/v1"

// IDPrefixes are allowed ID prefixes.
var IDPrefixes = []string{
	"user", "project", "pm", "role", "pi", "ptk", "pat", "container", "pipe_room", "deployment", "bearer",
	"ptkid", "secret", "gateway", "domain", "leap_token", "build",
}

// ValidateToken is used to validate a authentication token. Returns an error if the token is invalid.
func ValidateToken(authorization string) (string, error) {
	prefixSplit := strings.SplitN(authorization, "_", 2)
	if len(prefixSplit) != 2 {
		return "", types.InvalidToken("invalid authorization token format")
	}
	prefix := prefixSplit[0]
	for _, v := range IDPrefixes {
		if v == prefix {
			return prefix, nil
		}
	}
	return "", types.InvalidToken("invalid authorization token prefix: " + prefix)
}

// Client is used to define the Hop client. Please use the NewClient function to create this!
type Client struct {
	httpClient    *http.Client
	authorization string
	tokenType     string
	apiBase       string

	Pipe     *ClientCategoryPipe
	Projects *ClientCategoryProjects
	Ignite   *ClientCategoryIgnite
	Users    *ClientCategoryUsers
	Registry *ClientCategoryRegistry
	Channels *ClientCategoryChannels
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

// Does the specified HTTP request.
func (c *Client) do(ctx context.Context, a clientArgs) error {
	// Handle getting the body bytes.
	var r io.Reader
	if a.method != "GET" && a.body != nil {
		switch x := a.body.(type) {
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

	// Create the request.
	suffix := ""
	if a.query != nil {
		suffix = "?"
		first := false
		for k, v := range a.query {
			chunk := ""
			if first {
				first = false
			} else {
				chunk = "&"
			}
			suffix += chunk + url.QueryEscape(k) + "=" + url.QueryEscape(v)
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
	if r != nil {
		req.Header.Set("Content-Type", "application/json")
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

// Paginator is used to create a way to access paginated API routes.
type Paginator[T any] struct {
	c         *Client
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
func (p *Paginator[T]) Next(ctx context.Context) ([]T, error) {
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
	}); err != nil {
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
func (p *Paginator[T]) ForChunk(ctx context.Context, f func([]T) error) error {
	for a, err := p.Next(ctx); err != types.StopIteration; a, err = p.Next(ctx) {
		if err != nil {
			return err
		}
		if err = f(a); err != nil {
			return err
		}
	}
	return nil
}
