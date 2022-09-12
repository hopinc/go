package hop

import (
	"testing"

	"github.com/hopinc/hop-go/types"
)

func TestClient_Projects_Tokens_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "DELETE",
		wantPath:       "/projects/test%20test/tokens/test%20123",
		wantClientOpts: []ClientOption{WithProjectID("test test")},
		wantIgnore404:  false,
	}
	errorForTokenType(
		c, &ClientCategoryProjectsTokens{c: c},
		"Delete", []any{"test 123", WithProjectID("test test")}, "ptk")
	testApiSingleton(c,
		&ClientCategoryProjectsTokens{c: c},
		"Delete",
		[]any{"test 123", WithProjectID("test test")},
		nil)
}

func TestClient_Projects_Tokens_GetAll(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/projects/test%20test/tokens",
		wantResultKey:  "project_tokens",
		wantClientOpts: []ClientOption{WithProjectID("test test")},
		wantIgnore404:  false,
	}
	errorForTokenType(
		c, &ClientCategoryProjectsTokens{c: c},
		"GetAll", []any{WithProjectID("test test")}, "ptk")
	testApiSingleton(c,
		&ClientCategoryProjectsTokens{c: c},
		"GetAll",
		[]any{WithProjectID("test test")},
		[]*types.ProjectToken{{ID: "hello"}})
}

func TestClient_Projects_Tokens_Create(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "POST",
		wantPath:       "/projects/test%20test/tokens",
		wantResultKey:  "project_token",
		wantBody:       map[string][]types.ProjectPermission{"permissions": {}},
		wantClientOpts: []ClientOption{WithProjectID("test test")},
		wantIgnore404:  false,
	}
	errorForTokenType(
		c, &ClientCategoryProjectsTokens{c: c},
		"Create", []any{([]types.ProjectPermission)(nil), WithProjectID("test test")},
		"ptk")
	testApiSingleton(c,
		&ClientCategoryProjectsTokens{c: c},
		"Create",
		[]any{([]types.ProjectPermission)(nil), WithProjectID("test test")},
		&types.ProjectToken{ID: "hello"})
}

func TestClient_Projects_GetAllMembers(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/projects/test%20test/members",
		wantResultKey:  "members",
		wantClientOpts: []ClientOption{WithProjectID("test test")},
		wantIgnore404:  false,
	}
	testApiSingleton(c,
		&ClientCategoryProjects{c: c},
		"GetAllMembers",
		[]any{WithProjectID("test test")},
		[]*types.ProjectMember{{ID: "hello"}})
}

func TestClient_Projects_GetCurrentMember(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/projects/test/members/@me",
		wantResultKey:  "project_member",
		wantClientOpts: []ClientOption{WithProjectID("test")},
		wantIgnore404:  false,
	}
	errorForTokenType(c,
		&ClientCategoryProjects{c: c}, "GetCurrentMember",
		[]any{WithProjectID("test")}, "ptk")
	testApiSingleton(c,
		&ClientCategoryProjects{c: c},
		"GetCurrentMember",
		[]any{WithProjectID("test")},
		&types.ProjectMember{ID: "hello"})
}

func TestClient_Projects_Secrets_GetAll(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/projects/test%20test/secrets",
		wantResultKey:  "secrets",
		wantClientOpts: []ClientOption{WithProjectID("test test")},
		wantIgnore404:  false,
	}
	errorForTokenType(
		c, &ClientCategoryProjectsSecrets{c: c},
		"GetAll", []any{WithProjectID("test test")}, "ptk")
	testApiSingleton(c,
		&ClientCategoryProjectsSecrets{c: c},
		"GetAll",
		[]any{WithProjectID("test test")},
		[]*types.ProjectSecret{{ID: "hello"}})
}

func TestClient_Projects_Secrets_Create(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "PUT",
		wantPath:       "/projects/test%20test/secrets/ESCAPED%20SECRET",
		wantResultKey:  "secret",
		wantBody:       plainText("world"),
		wantClientOpts: []ClientOption{WithProjectID("test test")},
		wantIgnore404:  false,
	}
	errorForTokenType(
		c, &ClientCategoryProjectsSecrets{c: c},
		"Create", []any{"ESCAPED SECRET", "world", WithProjectID("test test")}, "ptk")
	testApiSingleton(c,
		&ClientCategoryProjectsSecrets{c: c},
		"Create",
		[]any{"ESCAPED SECRET", "world", WithProjectID("test test")},
		&types.ProjectSecret{ID: "hello"})
}

func TestClient_Projects_Secrets_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "DELETE",
		wantPath:       "/projects/test%20test/secrets/ESCAPED%20SECRET",
		wantClientOpts: []ClientOption{WithProjectID("test test")},
		wantIgnore404:  false,
	}
	errorForTokenType(
		c, &ClientCategoryProjectsSecrets{c: c},
		"Delete", []any{"ESCAPED SECRET", WithProjectID("test test")}, "ptk")
	testApiSingleton(c,
		&ClientCategoryProjectsSecrets{c: c},
		"Delete",
		[]any{"ESCAPED SECRET", WithProjectID("test test")},
		nil)
}
