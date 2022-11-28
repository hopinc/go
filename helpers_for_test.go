package hop

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockClientDoer struct {
	t *testing.T

	wantMethod     string
	wantPath       string
	wantQuery      map[string]string
	wantBody       any
	wantResultKey  string
	wantIgnore404  bool
	wantClientOpts []ClientOption
	tokenType      string

	// If you are using testApiSingleton, ignore these values.
	returnsResult any
	returnsErr    error
}

func (c *mockClientDoer) getProjectId(opts []ClientOption) string {
	projectId := ""
	for _, v := range opts {
		if v, ok := v.(projectIdOption); ok {
			projectId = v.projectId
		}
	}
	return projectId
}

func (c *mockClientDoer) getTokenType() string { return c.tokenType }

func (c *mockClientDoer) do(ctx context.Context, a ClientArgs, opts []ClientOption) error {
	c.t.Helper()
	assert.NotNil(c.t, ctx)
	assert.Equal(c.t, c.wantMethod, a.Method)
	assert.Equal(c.t, c.wantPath, a.Path)
	assert.Equal(c.t, c.wantQuery, a.Query)
	assert.Equal(c.t, c.wantBody, a.Body)
	assert.Equal(c.t, c.wantResultKey, a.ResultKey)
	assert.Equal(c.t, c.wantIgnore404, a.Ignore404)
	wantClientOpts := c.wantClientOpts
	if wantClientOpts == nil {
		wantClientOpts = []ClientOption{}
	}
	if opts == nil {
		opts = []ClientOption{}
	}
	assert.Equal(c.t, wantClientOpts, opts)
	if c.returnsErr != nil {
		return c.returnsErr
	}
	if c.returnsResult == nil {
		// Ensure the result is nil.
		assert.Nil(c.t, a.Result)
	} else {
		// Ensure that the result is a pointer to the type of the return result.
		pointerIsSame := assert.Equal(c.t, reflect.PtrTo(reflect.TypeOf(c.returnsResult)).String(), reflect.TypeOf(a.Result).String())

		if pointerIsSame {
			// Set the value of the result pointer to the return result.
			reflect.ValueOf(a.Result).Elem().Set(reflect.ValueOf(c.returnsResult))
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

func testApiSingleton(m *mockClientDoer, obj any, funcName string, args []any, result any) {
	t := m.t
	t.Helper()
	objReflect := reflect.Indirect(reflect.ValueOf(obj))
	funcResult := objReflect.MethodByName(funcName)
	if !funcResult.IsValid() {
		// The function doesn't exist.
		t.Fatalf("Function %s doesn't exist", funcName)
		return
	}

	// Check if this outputs 1 or 2 results. It should always end with a error.
	funcResultType := funcResult.Type()
	if funcResultType.NumOut() != 1 && funcResultType.NumOut() != 2 {
		t.Fatalf("Function %s should have 1 or 2 results", funcName)
		return
	}

	// Check the type of the last output is error.
	errType := reflect.TypeOf((*error)(nil)).Elem()
	if !funcResultType.Out(funcResultType.NumOut() - 1).Implements(errType) {
		t.Fatalf("Function %s should have an error as its last result", funcName)
		return
	}

	// Do 2 tests. One for client errors, one for no errors.
	runTest := func(testName string, err error) {
		t.Helper()
		t.Run(testName, func(t *testing.T) {
			m.returnsErr = err

			if result != nil {
				r := reflect.ValueOf(result)
				if result != nil && r.Type().Kind() == reflect.Ptr {
					// Dereference the result.
					r = r.Elem()
				}
				m.returnsResult = r.Interface()
			}

			reflectedArgs := make([]reflect.Value, len(args))
			for i, arg := range args {
				reflectedArgs[i] = reflect.ValueOf(arg)
			}
			results := funcResult.Call(append([]reflect.Value{reflect.ValueOf(context.Background())}, reflectedArgs...))
			var errValue reflect.Value
			if len(results) == 1 {
				errValue = results[0]
			} else {
				errValue = results[1]
			}
			assert.Equal(t, err, errValue.Interface())
			if err == nil && len(results) != 1 {
				// Check the result.
				assert.Equal(t, result, results[0].Interface())
			}
		})
	}
	runTest("client error", errors.New("cat tripped on wire"))
	runTest("no error", nil)
}

// Used to run a test to make sure it errors for token type.
func errorForTokenType(m *mockClientDoer, obj any, funcName string, args []any, tokenType string) {
	oldTokenType := m.tokenType
	defer func() { m.tokenType = oldTokenType }()
	m.tokenType = tokenType
	m.returnsErr = errors.New("shouldn't hit here")
	t := m.t
	t.Helper()
	objReflect := reflect.Indirect(reflect.ValueOf(obj))
	funcResult := objReflect.MethodByName(funcName)
	if !funcResult.IsValid() {
		// The function doesn't exist.
		t.Fatalf("Function %s doesn't exist", funcName)
		return
	}
	funcResultType := funcResult.Type()
	if funcResultType.NumOut() == 0 || funcResultType.Out(funcResultType.NumOut()-1).String() != "error" {
		t.Fatalf("Function %s should have an error as its last result", funcName)
		return
	}

	t.Run("bad token type", func(t *testing.T) {
		reflectedArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			reflectedArgs[i] = reflect.ValueOf(arg)
		}
		results := funcResult.Call(append([]reflect.Value{reflect.ValueOf(context.Background())}, reflectedArgs...))
		errValue := results[funcResultType.NumOut()-1]
		assert.NotNil(t, errValue.Interface())
	})
}
