package hop

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.hop.io/sdk/types"
)

type customError struct{}

func (c customError) Error() string { return "capybara nibbled wire" }

type errorReader struct {
	err error
}

func (e errorReader) Read([]byte) (n int, err error) {
	return 0, e.err
}

func Test_handleErrors(t *testing.T) {
	tests := []struct {
		name string

		status int
		body   []byte // if nil, will error with "capybara nibbled wire" and customError.

		expectsErr     string
		expectsErrType any
	}{
		{
			name:           "body error",
			status:         400,
			body:           nil,
			expectsErr:     "capybara nibbled wire",
			expectsErrType: customError{},
		},
		{
			name:   "unable to unmarshal",
			status: 400,
			body:   []byte("testing testing 123"),
			expectsErr: "status code 400 (cannot unmarshal from json): " +
				"testing testing 123",
			expectsErrType: types.ServerError(""),
		},
		{
			name:           "invalid auth",
			status:         403,
			body:           []byte(`{"success":false,"error":{"code":"invalid_auth","message":"invalid auth"}}`),
			expectsErr:     "invalid auth",
			expectsErrType: types.NotAuthorized(""),
		},
		{
			name:           "bad request",
			status:         400,
			body:           []byte(`{"success":false,"error":{"code":"bad_request","message":"oof"}}`),
			expectsErr:     "bad_request: oof",
			expectsErrType: types.BadRequest{},
		},
		{
			name:           "not found",
			status:         404,
			body:           []byte(`{"success":false,"error":{"code":"not_found","message":"oof"}}`),
			expectsErr:     "oof",
			expectsErrType: types.NotFound{},
		},
		{
			name:           "server error",
			status:         500,
			body:           []byte(`{"success":false,"error":{"code":"server_error","message":"oof"}}`),
			expectsErr:     "oof",
			expectsErrType: types.ServerError(""),
		},
		{
			name:           "unknown server error",
			status:         401,
			body:           []byte(`{"success":false,"error":{"code":"unknown_error","message":"oof"}}`),
			expectsErr:     "status code 401 (unknown_error): oof",
			expectsErrType: types.UnknownServerError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				StatusCode: tt.status,
				Body:       io.NopCloser(bytes.NewReader(tt.body)),
			}
			if tt.body == nil {
				res.Body = io.NopCloser(errorReader{customError{}})
			}

			err := handleErrors(res)
			if tt.expectsErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectsErr)
				assert.IsType(t, tt.expectsErrType, err)
			}
		})
	}
}
