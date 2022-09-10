package hop

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockClientDoer struct {
	t *testing.T

	wantMethod    string
	wantPath      string
	wantQuery     map[string]string
	wantBody      any
	wantResultKey string
	wantIgnore404 bool
	tokenType     string

	returnsResult any
	returnsErr    error
}

func (c *mockClientDoer) getTokenType() string { return c.tokenType }

func (c *mockClientDoer) do(ctx context.Context, a clientArgs) error {
	c.t.Helper()
	assert.NotNil(c.t, ctx)
	assert.Equal(c.t, c.wantMethod, a.method)
	assert.Equal(c.t, c.wantPath, a.path)
	assert.Equal(c.t, c.wantQuery, a.query)
	assert.Equal(c.t, c.wantBody, a.body)
	assert.Equal(c.t, c.wantResultKey, a.resultKey)
	assert.Equal(c.t, c.wantIgnore404, a.ignore404)
	if c.returnsErr != nil {
		return c.returnsErr
	}
	if c.returnsResult == nil {
		// Ensure the result is nil.
		assert.Nil(c.t, a.result)
	} else {
		// Ensure that the result is a pointer to the type of the return result.
		if assert.Equal(c.t, reflect.PtrTo(reflect.TypeOf(c.returnsResult)), reflect.TypeOf(a.result)) {
			// Set the value of the result pointer to the return result.
			reflect.ValueOf(a.result).Elem().Set(reflect.ValueOf(c.returnsResult))
		}
	}
	return nil
}

func dup[T any](x T, l int) []T {
	r := make([]T, l)
	for i := range r {
		r[i] = x
	}
	return r
}

func rawify(m map[string]any) map[string]json.RawMessage {
	r := make(map[string]json.RawMessage, len(m))
	for k, v := range m {
		r[k], _ = json.Marshal(v)
	}
	return r
}
