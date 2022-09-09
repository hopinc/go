package hop

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hopinc/hop-go/types"
	"github.com/stretchr/testify/assert"
)

func TestClient_SetAPIBase(t *testing.T) {
	tests := []struct {
		name string

		baseUrl string
		expects string
	}{
		{
			name:    "remove trailing slash",
			baseUrl: "https://example.com/",
			expects: "https://example.com",
		},
		{
			name:    "add https",
			baseUrl: "example.com",
			expects: "https://example.com",
		},
		{
			name:    "allow http",
			baseUrl: "http://example.com",
			expects: "http://example.com",
		},
		{
			name:    "full url",
			baseUrl: "https://example.com/v1",
			expects: "https://example.com/v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			c.SetAPIBase(tt.baseUrl)
			assert.Equal(t, tt.expects, c.apiBase)
		})
	}
}

type mockHttpRoundTripper struct {
	t *testing.T

	err error

	wantHeaders http.Header
	wantUrl     string
	wantBody    string

	returnsStatus int
	returnsBody   string
}

func (h mockHttpRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	h.t.Helper()

	assert.Equal(h.t, h.wantHeaders, req.Header)
	assert.Equal(h.t, h.wantUrl, req.URL.String())
	if req.Body == nil {
		assert.Equal(h.t, "", h.wantBody)
	} else {
		body, err := io.ReadAll(req.Body)
		assert.NoError(h.t, err)
		assert.Equal(h.t, h.wantBody, string(body))
	}

	if h.err != nil {
		return nil, h.err
	}

	return &http.Response{
		StatusCode: h.returnsStatus,
		Body:       io.NopCloser(strings.NewReader(h.returnsBody)),
	}, nil
}

func makeJsonSyntaxError(offset int) error {
	x := make([]byte, offset)
	for i := range x {
		x[i] = '1'
	}
	x[offset-1] = '{'
	return json.Unmarshal(x, &struct{}{})
}

type marshalFail struct{}

func (marshalFail) MarshalJSON() ([]byte, error) {
	return nil, errors.New("marshal fail")
}

func TestClient_do(t *testing.T) {
	tests := []struct {
		name string

		wantHeaders http.Header
		wantUrl     string
		wantBody    string

		returnsStatus int
		returnsBody   string
		returnsError  error

		baseUrl      string
		method       string
		path         string
		resultKey    string
		query        map[string]string
		body         any
		ignore404    bool
		expectsError error
	}{
		{
			name: "http client error",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:      "https://api.hop.io/v1/test",
			returnsError: errors.New("hamster tripped on wire"),
			expectsError: errors.New("hamster tripped on wire"),
			method:       "GET",
			path:         "/test",
		},
		{
			name: "query params",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:       "https://api.hop.io/v1/test?foo=bar",
			returnsBody:   `{"data":{"foo":"bar"}}`,
			returnsStatus: 200,
			method:        "GET",
			path:          "/test",
			query: map[string]string{
				"foo": "bar",
			},
		},
		{
			name: "multiple query params",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:       "https://api.hop.io/v1/test?foo=bar&foo2=bar2",
			returnsBody:   `{"data":{"foo":"bar"}}`,
			returnsStatus: 200,
			method:        "GET",
			path:          "/test",
			query: map[string]string{
				"foo":  "bar",
				"foo2": "bar2",
			},
		},
		{
			name: "get success",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			returnsBody:   `{"data":{"foo":"bar"}}`,
			returnsStatus: 200,
			method:        "GET",
			path:          "/test",
		},
		{
			name: "no content",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			returnsStatus: 204,
			method:        "GET",
			path:          "/test",
		},
		{
			name: "get ignore 404",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			returnsBody:   `{"data":{"foo":"bar"}}`,
			returnsStatus: 404,
			method:        "GET",
			path:          "/test",
			ignore404:     true,
		},
		{
			name: "get ignore body",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			returnsBody:   `{"data":{"foo":"bar"}}`,
			returnsStatus: 200,
			method:        "GET",
			path:          "/test",
			body:          map[string]string{},
		},
		{
			name: "pluck key",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			returnsBody:   `{"data":{"pog":{"foo":"bar"}}}`,
			returnsStatus: 200,
			resultKey:     "pog",
			method:        "GET",
			path:          "/test",
			body:          map[string]string{},
		},
		{
			name: "bad response json",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			returnsBody:   "{",
			returnsStatus: 200,
			expectsError:  makeJsonSyntaxError(1),
			resultKey:     "pog",
			method:        "GET",
			path:          "/test",
		},
		{
			name: "marshalled json post request",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
				"Content-Type":  {"application/json"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			wantBody:      `{"hello":"world"}`,
			returnsBody:   `{"data":{"foo":"bar"}}`,
			returnsStatus: 200,
			method:        "POST",
			path:          "/test",
			body:          map[string]string{"hello": "world"},
		},
		{
			name: "json bytes post request",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
				"Content-Type":  {"application/json"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			wantBody:      `{"hello":"world"}`,
			returnsBody:   `{"data":{"foo":"bar"}}`,
			returnsStatus: 200,
			method:        "POST",
			path:          "/test",
			body:          []byte(`{"hello":"world"}`),
		},
		{
			name: "plain text post request",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
				"Content-Type":  {"text/plain"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			wantBody:      "hello world",
			returnsBody:   `{"data":{"foo":"bar"}}`,
			returnsStatus: 200,
			method:        "POST",
			path:          "/test",
			body:          plainText("hello world"),
		},
		{
			name:         "body marshal error",
			expectsError: errors.New("marshal fail"),
			method:       "POST",
			path:         "/test",
			body:         marshalFail{},
		},
		{
			name: "calls handleErrors",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			returnsBody:   `{"error":{"message":"fail","code":"not_found"}}`,
			expectsError:  types.NotFound("fail"),
			returnsStatus: 404,
			method:        "GET",
			path:          "/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient: &http.Client{
					Transport: mockHttpRoundTripper{
						t:             t,
						err:           tt.returnsError,
						wantHeaders:   tt.wantHeaders,
						wantUrl:       tt.wantUrl,
						wantBody:      tt.wantBody,
						returnsStatus: tt.returnsStatus,
						returnsBody:   tt.returnsBody,
					},
				},
				authorization: "testing",
				isTest:        true,
			}
			if tt.baseUrl != "" {
				c.SetAPIBase(tt.baseUrl)
			}
			result := map[string]string{}
			var ptr any = &result
			if tt.returnsBody == "" {
				// blank responses should be a nil pointer
				ptr = nil
			}
			err := c.do(context.TODO(), clientArgs{
				method:    tt.method,
				path:      tt.path,
				resultKey: tt.resultKey,
				query:     tt.query,
				body:      tt.body,
				result:    ptr,
				ignore404: tt.ignore404,
			})
			if tt.expectsError == nil {
				assert.NoError(t, err)
			} else {
				if err2 := errors.Unwrap(err); err2 != nil {
					err = err2
				}
				assert.Equal(t, tt.expectsError, err)
				return
			}
			if ptr != nil {
				// Check the body is what we expect.
				assert.Equal(t, map[string]string{"foo": "bar"}, result)
			}
		})
	}
}
