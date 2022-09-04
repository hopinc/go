package hopgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// APIBase is used to define the base API URL.
const APIBase = "https://api.hop.io/v1"

// Client is used to define the Hop client. Please use the NewClient function to create this!
type Client struct {
	httpClient    *http.Client
	authorization string

	Pipe     *ClientCategoryPipe
	Projects *ClientCategoryProjects
	Ignite   *ClientCategoryIgnite
	Users    *ClientCategoryUsers
	Registry *ClientCategoryRegistry
	Channels *ClientCategoryChannels
}

// NewClient is used to make a new Hop client.
func NewClient(authorization string) *Client {
	var c Client
	c = Client{
		httpClient:    &http.Client{},
		authorization: authorization,
		Pipe:          newPipe(&c),
		Projects:      newProjects(&c),
		Ignite:        newIgnite(&c),
		Users:         newUsers(&c),
		Registry:      newRegistry(&c),
		Channels:      newChannels(&c),
	}
	return &c
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
	req, err := http.NewRequestWithContext(ctx, a.method, APIBase+a.path+suffix, r)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.authorization)
	req.Header.Set("Accept", "application/json")

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
