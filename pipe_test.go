package hop

import (
	"testing"

	"github.com/hopinc/hop-go/types"
)

func TestClient_Pipe_Rooms_GetAll(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/pipe/rooms",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantResultKey:  "rooms",
		wantIgnore404:  false,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryPipeRooms{c: c},
		"GetAll",
		[]any{WithProjectID("test123")},
		[]*types.Room{{Name: "hello"}})
}

func TestClient_Pipe_Rooms_Create(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "POST",
		wantPath:      "/pipe/rooms",
		wantResultKey: "room",
		wantIgnore404: false,
		wantBody:      types.RoomCreationOptions{Name: "test"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryPipeRooms{c: c},
		"Create",
		[]any{types.RoomCreationOptions{Name: "test"}},
		&types.Room{Name: "hello"})
}

func TestClient_Pipe_Rooms_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "DELETE",
		wantPath:       "/pipe/rooms/test%20test",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantIgnore404:  false,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryPipeRooms{c: c},
		"Delete",
		[]any{"test test", WithProjectID("test123")},
		nil)
}
