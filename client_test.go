package hop

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hopinc/hop-go/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name string

		token         string
		expectsPrefix string
		expectsErr    string
	}{
		{
			name:          "user",
			token:         "user_xyz",
			expectsPrefix: "user",
		},
		{
			name:          "project",
			token:         "project_xyz",
			expectsPrefix: "project",
		},
		{
			name:          "pm",
			token:         "pm_xyz",
			expectsPrefix: "pm",
		},
		{
			name:          "role",
			token:         "role_xyz",
			expectsPrefix: "role",
		},
		{
			name:          "pi",
			token:         "pi_xyz",
			expectsPrefix: "pi",
		},
		{
			name:          "ptk",
			token:         "ptk_xyz",
			expectsPrefix: "ptk",
		},
		{
			name:          "pat",
			token:         "pat_xyz",
			expectsPrefix: "pat",
		},
		{
			name:          "container",
			token:         "container_xyz",
			expectsPrefix: "container",
		},
		{
			name:          "pipe_room",
			token:         "pipe_room_xyz",
			expectsPrefix: "pipe_room",
		},
		{
			name:          "deployment",
			token:         "deployment_xyz",
			expectsPrefix: "deployment",
		},
		{
			name:          "bearer",
			token:         "bearer_xyz",
			expectsPrefix: "bearer",
		},
		{
			name:          "ptkid",
			token:         "ptkid_xyz",
			expectsPrefix: "ptkid",
		},
		{
			name:          "secret",
			token:         "secret_xyz",
			expectsPrefix: "secret",
		},
		{
			name:          "gateway",
			token:         "gateway_xyz",
			expectsPrefix: "gateway",
		},
		{
			name:          "domain",
			token:         "domain_xyz",
			expectsPrefix: "domain",
		},
		{
			name:          "leap_token",
			token:         "leap_token_xyz",
			expectsPrefix: "leap_token",
		},
		{
			name:          "build",
			token:         "build_xyz",
			expectsPrefix: "build",
		},
		{
			name:       "invalid",
			token:      "invalid_xyz",
			expectsErr: "invalid authorization token prefix: invalid_xyz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, err := ValidateToken(tt.token)
			if tt.expectsErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectsErr)
			}
			assert.Equal(t, tt.expectsPrefix, prefix)
		})
	}
}

func noNilPointers(x any) bool {
	v := reflect.Indirect(reflect.ValueOf(x))
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanSet() && field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return false
			}
			if !noNilPointers(field.Interface()) {
				return false
			}
		}
	}
	return true
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string

		token      string
		tokenType  string
		expectsErr string
	}{
		{
			name:      "bearer",
			token:     "bearer_xyz",
			tokenType: "bearer",
		},
		{
			name:      "pat",
			token:     "pat_xyz",
			tokenType: "pat",
		},
		{
			name:      "ptk",
			token:     "ptk_xyz",
			tokenType: "ptk",
		},
		{
			name:       "invalid",
			token:      "invalid_xyz",
			expectsErr: "invalid authorization token prefix: invalid_xyz",
		},
		{
			name:       "disallowed token type",
			token:      "user_xyz",
			expectsErr: "invalid authorization token prefix: user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.token)
			if tt.expectsErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.tokenType, c.getTokenType())
				assert.True(t, noNilPointers(c))
			} else {
				assert.EqualError(t, err, tt.expectsErr)
			}
		})
	}
}

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

func TestClient_AddClientOptions(t *testing.T) {
	c := &Client{}
	c.AddClientOptions(WithProjectID("a"))
	c.AddClientOptions(WithProjectID("b"), WithProjectID("c"))
	assert.Equal(t, []ClientOption{projectIdOption{projectId: "a"}, projectIdOption{projectId: "b"},
		projectIdOption{projectId: "c"}}, c.opts)
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
		clientOpts   []ClientOption
		funcOpts     []ClientOption
		expectsError error
	}{
		{
			name: "http client error",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
			},
			wantUrl:       "https://api.hop.io/v1/test",
			returnsStatus: 204,
			method:        "GET",
			path:          "/test",
		},
		{
			name: "client project id option applied",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
				"User-Agent":    {userAgent},
			},
			wantUrl:       "https://api.hop.io/v1/test?project=test",
			returnsStatus: 204,
			method:        "GET",
			path:          "/test",
			clientOpts:    []ClientOption{projectIdOption{projectId: "test"}},
		},
		{
			name: "function project id option applied",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
				"User-Agent":    {userAgent},
			},
			wantUrl:       "https://api.hop.io/v1/test?project=test",
			returnsStatus: 204,
			method:        "GET",
			path:          "/test",
			funcOpts:      []ClientOption{projectIdOption{projectId: "test"}},
		},
		{
			name: "function project id option overrides client option",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
				"User-Agent":    {userAgent},
			},
			wantUrl:       "https://api.hop.io/v1/test?project=test2",
			returnsStatus: 204,
			method:        "GET",
			path:          "/test",
			clientOpts:    []ClientOption{projectIdOption{projectId: "test1"}},
			funcOpts:      []ClientOption{projectIdOption{projectId: "test2"}},
		},
		{
			name: "custom base url",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
				"User-Agent":    {userAgent},
			},
			baseUrl:       "https://example.com/v9999",
			wantUrl:       "https://example.com/v9999/test",
			returnsStatus: 204,
			method:        "GET",
			path:          "/test",
		},
		{
			name: "get ignore 404",
			wantHeaders: http.Header{
				"Accept":        {"application/json"},
				"Authorization": {"testing"},
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				"User-Agent":    {userAgent},
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
				opts:          tt.clientOpts,
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
			err := c.do(context.Background(), clientArgs{
				method:    tt.method,
				path:      tt.path,
				resultKey: tt.resultKey,
				query:     tt.query,
				body:      tt.body,
				result:    ptr,
				ignore404: tt.ignore404,
			}, tt.funcOpts)
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

func TestPaginator_Next(t *testing.T) {
	tests := []struct {
		name string

		queries         []map[string]string
		clientResults   []map[string]json.RawMessage
		clientErrors    []error
		returnedResults [][]string
		returnedErrors  []error

		offsetStrat  bool
		initialQuery map[string]string
	}{
		{
			name:            "first client error",
			queries:         dup(map[string]string{"orderBy": "asc", "page": "1", "sortBy": "test_sort"}, 1),
			clientResults:   dup((map[string]json.RawMessage)(nil), 1),
			clientErrors:    []error{errors.New("fail")},
			returnedResults: dup([]string(nil), 1),
			returnedErrors:  []error{errors.New("fail")},
			offsetStrat:     false,
			initialQuery:    nil,
		},
		{
			name: "second client error",
			queries: []map[string]string{
				{"orderBy": "asc", "page": "1", "sortBy": "test_sort"},
				{"orderBy": "asc", "page": "2", "sortBy": "test_sort"},
			},
			clientResults: []map[string]json.RawMessage{rawify(map[string]any{
				"total_count": 4,
				"items":       []string{"a", "b"},
			}), nil},
			clientErrors:    []error{nil, errors.New("fail")},
			returnedResults: [][]string{{"a", "b"}, nil},
			returnedErrors:  []error{nil, errors.New("fail")},
			offsetStrat:     false,
			initialQuery:    nil,
		},
		{
			name: "query param passthrough",
			queries: []map[string]string{
				{"orderBy": "asc", "page": "1", "sortBy": "test_sort", "x": "y"},
			},
			clientResults: []map[string]json.RawMessage{rawify(map[string]any{
				"total_count": 2,
				"items":       []string{"a", "b"},
			}), nil},
			clientErrors:    []error{nil},
			returnedResults: [][]string{{"a", "b"}},
			returnedErrors:  []error{nil},
			offsetStrat:     false,
			initialQuery:    map[string]string{"x": "y"},
		},
		{
			name: "page mode single success",
			queries: []map[string]string{
				{"orderBy": "asc", "page": "1", "sortBy": "test_sort"},
			},
			clientResults: []map[string]json.RawMessage{rawify(map[string]any{
				"total_count": 2,
				"items":       []string{"a", "b"},
			}), nil},
			clientErrors:    []error{nil},
			returnedResults: [][]string{{"a", "b"}},
			returnedErrors:  []error{nil},
			offsetStrat:     false,
		},
		{
			name: "page mode multiple success",
			queries: []map[string]string{
				{"orderBy": "asc", "page": "1", "sortBy": "test_sort"},
				{"orderBy": "asc", "page": "2", "sortBy": "test_sort"},
			},
			clientResults: []map[string]json.RawMessage{rawify(map[string]any{
				"total_count": 4,
				"items":       []string{"a", "b"},
			}), rawify(map[string]any{
				"total_count": 4,
				"items":       []string{"c", "d"},
			})},
			clientErrors:    []error{nil, nil},
			returnedResults: [][]string{{"a", "b"}, {"c", "d"}},
			returnedErrors:  []error{nil, nil},
			offsetStrat:     false,
		},
		{
			name: "page mode stop iteration client hit",
			queries: []map[string]string{
				{"orderBy": "asc", "page": "1", "sortBy": "test_sort"},
				{},
			},
			clientResults: []map[string]json.RawMessage{rawify(map[string]any{
				"total_count": 2,
				"items":       []string{"a", "b"},
			}), nil},
			clientErrors:    []error{nil, nil},
			returnedResults: [][]string{{"a", "b"}, {}},
			returnedErrors:  []error{nil, types.StopIteration},
			offsetStrat:     false,
		},
		{
			name: "page mode stop iteration api hit",
			queries: []map[string]string{
				{"orderBy": "asc", "page": "1", "sortBy": "test_sort"},
				{"orderBy": "asc", "page": "2", "sortBy": "test_sort"},
			},
			clientResults: []map[string]json.RawMessage{rawify(map[string]any{
				"total_count": 4,
				"items":       []string{"a", "b"},
			}), rawify(map[string]any{
				"total_count": 4,
				"items":       []string{},
			})},
			clientErrors:    []error{nil, nil},
			returnedResults: [][]string{{"a", "b"}, {}},
			returnedErrors:  []error{nil, types.StopIteration},
			offsetStrat:     false,
		},
		{
			name: "offset mode single success",
			queries: []map[string]string{
				{"limit": "2", "offset": "0", "orderBy": "asc", "sortBy": "test_sort"},
			},
			clientResults: []map[string]json.RawMessage{rawify(map[string]any{
				"items": []string{"a", "b"},
			}), nil},
			clientErrors:    []error{nil},
			returnedResults: [][]string{{"a", "b"}},
			returnedErrors:  []error{nil},
			offsetStrat:     true,
		},
		{
			name: "offset mode multiple success",
			queries: []map[string]string{
				{"limit": "2", "offset": "0", "orderBy": "asc", "sortBy": "test_sort"},
				{"limit": "2", "offset": "2", "orderBy": "asc", "sortBy": "test_sort"},
			},
			clientResults: []map[string]json.RawMessage{rawify(map[string]any{
				"items": []string{"a", "b"},
			}), rawify(map[string]any{
				"items": []string{"c", "d"},
			})},
			clientErrors:    []error{nil, nil},
			returnedResults: [][]string{{"a", "b"}, {"c", "d"}},
			returnedErrors:  []error{nil, nil},
			offsetStrat:     true,
		},
		{
			name: "offset mode stop iteration",
			queries: []map[string]string{
				{"limit": "2", "offset": "0", "orderBy": "asc", "sortBy": "test_sort"},
				{"limit": "2", "offset": "2", "orderBy": "asc", "sortBy": "test_sort"},
			},
			clientResults: []map[string]json.RawMessage{rawify(map[string]any{
				"items": []string{"a", "b"},
			}), rawify(map[string]any{
				"items": []string{},
			})},
			clientErrors:    []error{nil, nil},
			returnedResults: [][]string{{"a", "b"}, {"c", "d"}},
			returnedErrors:  []error{nil, types.StopIteration},
			offsetStrat:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &mockClientDoer{
				t:             t,
				wantMethod:    "GET",
				wantPath:      "/test",
				wantQuery:     tt.queries[0],
				returnsResult: tt.clientResults[0],
				returnsErr:    tt.clientErrors[0],
			}
			p := &Paginator[string]{
				c:           d,
				total:       -1,
				path:        "/test",
				resultKey:   "items",
				sortBy:      "test_sort",
				query:       tt.initialQuery,
				limit:       2,
				offsetStrat: tt.offsetStrat,
			}
			for i := 0; i < len(tt.queries); i++ {
				d.wantQuery = tt.queries[i]
				d.returnsResult = tt.clientResults[i]
				d.returnsErr = tt.clientErrors[i]
				res, err := p.Next(context.Background())
				if tt.returnedErrors[i] == nil {
					// Make sure error is nil.
					assert.NoError(t, err)

					// Make sure the result is what we expect.
					assert.Equal(t, tt.returnedResults[i], res)
				} else {
					assert.Equal(t, tt.returnedErrors[i], err)
				}
			}
		})
	}
}

func TestPaginator_ForChunk(t *testing.T) {
	t.Run("run until stop", func(t *testing.T) {
		clientDoers := []*mockClientDoer{
			{
				t:             t,
				wantMethod:    "GET",
				wantPath:      "/test",
				wantQuery:     map[string]string{"orderBy": "asc", "page": "1", "sortBy": "test_sort"},
				returnsResult: rawify(map[string]any{"total_count": 4, "items": []string{"a", "b"}}),
				returnsErr:    nil,
			},
			{
				t:             t,
				wantMethod:    "GET",
				wantPath:      "/test",
				wantQuery:     map[string]string{"orderBy": "asc", "page": "2", "sortBy": "test_sort"},
				returnsResult: rawify(map[string]any{"total_count": 4, "items": []string{"c"}}),
				returnsErr:    nil,
			},
			{
				t:             t,
				wantMethod:    "GET",
				wantPath:      "/test",
				wantQuery:     map[string]string{"orderBy": "asc", "page": "3", "sortBy": "test_sort"},
				returnsResult: rawify(map[string]any{"total_count": 4, "items": []string{"d"}}),
				returnsErr:    nil,
			},
			{
				t:             t,
				wantMethod:    "GET",
				wantPath:      "/test",
				wantQuery:     map[string]string{"orderBy": "asc", "page": "4", "sortBy": "test_sort"},
				returnsResult: rawify(map[string]any{"total_count": 4, "items": []string{}}),
				returnsErr:    types.StopIteration,
			},
		}
		s := []string{}
		p := &Paginator[string]{
			total:     -1,
			path:      "/test",
			resultKey: "items",
			sortBy:    "test_sort",
		}
		p.c = clientDoers[0]
		i := 1
		err := p.ForChunk(context.Background(), func(res []string) error {
			s = append(s, res...)
			p.c = clientDoers[i]
			i++
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c", "d"}, s)
	})

	t.Run("loop error", func(t *testing.T) {
		p := &Paginator[string]{
			c: &mockClientDoer{
				t:             t,
				wantMethod:    "GET",
				wantPath:      "/test",
				wantQuery:     map[string]string{"orderBy": "asc", "page": "1", "sortBy": "test_sort"},
				returnsResult: rawify(map[string]any{"total_count": 4, "items": []string{"a", "b"}}),
				returnsErr:    nil,
			},
			total:     -1,
			path:      "/test",
			resultKey: "items",
			sortBy:    "test_sort",
		}
		err := p.ForChunk(context.Background(), func(res []string) error {
			return errors.New("test error")
		})
		assert.EqualError(t, err, "test error")
	})

	t.Run("client error", func(t *testing.T) {
		clientDoers := []*mockClientDoer{
			{
				t:             t,
				wantMethod:    "GET",
				wantPath:      "/test",
				wantQuery:     map[string]string{"orderBy": "asc", "page": "1", "sortBy": "test_sort"},
				returnsResult: rawify(map[string]any{"total_count": 4, "items": []string{"a", "b"}}),
				returnsErr:    errors.New("test error"),
			},
		}
		p := &Paginator[string]{
			c:         clientDoers[0],
			total:     -1,
			path:      "/test",
			resultKey: "items",
			sortBy:    "test_sort",
		}
		err := p.ForChunk(context.Background(), func(res []string) error {
			return nil
		})
		assert.EqualError(t, err, "test error")
	})
}
