package hop

import (
	"testing"

	"github.com/hopinc/hop-go/types"
)

func TestClient_Users_Me_Get(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/users/@me",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	errorForTokenType(c, &ClientCategoryUsersMe{c: c}, "Get", []any{}, "ptk")
	testApiSingleton(c,
		&ClientCategoryUsersMe{c: c},
		"Get",
		[]any{},
		&types.UserMeInfo{User: types.User{Name: "jeff"}})
}

func TestClient_Users_Me_CreatePat(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "POST",
		wantPath:      "/users/@me/pats",
		wantIgnore404: false,
		wantBody:      map[string]string{"name": "test"},
		wantResultKey: "pat",
		tokenType:     "pat",
	}
	errorForTokenType(c, &ClientCategoryUsersMe{c: c}, "CreatePat", []any{"test"}, "ptk")
	testApiSingleton(c,
		&ClientCategoryUsersMe{c: c},
		"CreatePat",
		[]any{"test"},
		&types.UserPat{PAT: "hello"})
}

func TestClient_Users_Me_GetAllPats(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/users/@me/pats",
		wantIgnore404: false,
		wantResultKey: "pats",
		tokenType:     "pat",
	}
	errorForTokenType(c, &ClientCategoryUsersMe{c: c}, "GetAllPats", []any{}, "ptk")
	testApiSingleton(c,
		&ClientCategoryUsersMe{c: c},
		"GetAllPats",
		[]any{},
		[]*types.UserPat{{PAT: "hello"}})
}

func TestClient_Users_Me_DeletePat(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "DELETE",
		wantPath:      "/users/@me/pats/test%20test",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	errorForTokenType(c, &ClientCategoryUsersMe{c: c}, "DeletePat", []any{"test test"}, "ptk")
	testApiSingleton(c,
		&ClientCategoryUsersMe{c: c},
		"DeletePat",
		[]any{"test test"},
		nil)
}
