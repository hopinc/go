package hop

import (
	"testing"

	"github.com/hopinc/hop-go/types"
	"github.com/stretchr/testify/assert"
)

func TestClient_Channels_Create(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "POST",
		wantPath:      "/channels",
		wantResultKey: "channel",
		wantIgnore404: false,
		wantQuery:     map[string]string{"project": "test123"},
		wantBody:      map[string]any{"type": types.ChannelTypePrivate, "state": map[string]any{"do_you_like_waffles": true}},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"Create",
		[]any{types.ChannelTypePrivate, map[string]any{"do_you_like_waffles": true}, "", "test123"},
		&types.Channel{ID: "hello"})
}

func TestClient_Channels_Get(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/channels/test%20test",
		wantQuery:     map[string]string{"project": "test123"},
		wantResultKey: "channel",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"Get",
		[]any{"test123", "test test"},
		&types.Channel{ID: "hello"})
}

func TestClient_Channels_GetAll(t *testing.T) {
	c := &mockClientDoer{}
	res := (&ClientCategoryChannels{c: c}).GetAll(
		"test123")
	assert.Equal(t, res, &Paginator[*types.Channel]{
		c:         c,
		total:     -1,
		path:      "/channels",
		resultKey: "channels",
		sortBy:    "created_at",
		query:     map[string]string{"project": "test123"},
	})
}

func TestClient_Channels_SubscribeToken(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "PUT",
		wantPath:      "/channels/test%20test/subscribers/testing%20testing123",
		wantQuery:     map[string]string{"project": "test123"},
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"SubscribeToken",
		[]any{"test123", "test test", "testing testing123"},
		nil)
}

func TestClient_Channels_SetState(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "PUT",
		wantPath:      "/channels/test%20test/state",
		wantQuery:     map[string]string{"project": "test123"},
		wantIgnore404: false,
		wantBody:      map[string]any{"do_you_like_waffles": true},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"SetState",
		[]any{"test123", "test test", map[string]any{"do_you_like_waffles": true}},
		nil)
}

func TestClient_Channels_PatchState(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "PATCH",
		wantPath:      "/channels/test%20test/state",
		wantQuery:     map[string]string{"project": "test123"},
		wantIgnore404: false,
		wantBody:      map[string]any{"do_you_like_waffles": true},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"PatchState",
		[]any{"test123", "test test", map[string]any{"do_you_like_waffles": true}},
		nil)
}

func TestClient_Channels_PublishMessage(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "POST",
		wantPath:      "/channels/test%20test/messages",
		wantQuery:     map[string]string{"project": "test123"},
		wantIgnore404: false,
		wantBody:      map[string]any{"e": "hello", "d": "world"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"PublishMessage",
		[]any{"test123", "test test", "hello", "world"},
		nil)
}

func TestClient_Channels_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "DELETE",
		wantPath:      "/channels/test%20test",
		wantQuery:     map[string]string{"project": "test123"},
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"Delete",
		[]any{"test123", "test test"},
		nil)
}

func TestClient_Channels_GetStats(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/channels/test%20test/stats",
		wantQuery:     map[string]string{"project": "test123"},
		wantResultKey: "stats",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"GetStats",
		[]any{"test123", "test test"},
		&types.Stats{OnlineCount: 9000})
}

func TestClient_Channels_Tokens_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "DELETE",
		wantPath:      "/channels/tokens/test%20test123",
		wantQuery:     map[string]string{"project": "test123"},
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"Delete",
		[]any{"test123", "test test123"},
		nil)
}

func TestClient_Channels_Tokens_Create(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "POST",
		wantPath:      "/channels/tokens",
		wantQuery:     map[string]string{"project": "test123"},
		wantResultKey: "token",
		wantIgnore404: false,
		wantBody:      map[string]any{"state": map[string]any{"hello": "world"}},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"Create",
		[]any{"test123", map[string]any{"hello": "world"}},
		&types.ChannelToken{ID: "hello"})
}

func TestClient_Channels_Tokens_SetState(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "PATCH",
		wantPath:      "/channels/tokens/test%20test123",
		wantQuery:     map[string]string{"project": "test123"},
		wantIgnore404: false,
		wantBody:      map[string]any{"state": map[string]any{"do_you_like_waffles": true}},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"SetState",
		[]any{"test123", "test test123", map[string]any{"do_you_like_waffles": true}},
		nil)
}

func TestClient_Channels_Tokens_Get(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/channels/tokens/test%20test123",
		wantQuery:     map[string]string{"project": "test123"},
		wantResultKey: "token",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"Get",
		[]any{"test123", "test test123"},
		&types.ChannelToken{ID: "hello"})
}

func TestClient_Channels_Tokens_PublishDirectMessage(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "POST",
		wantPath:      "/channels/tokens/test%20test123/messages",
		wantQuery:     map[string]string{"project": "test123"},
		wantIgnore404: false,
		wantBody:      map[string]any{"e": "hello", "d": "world"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"PublishDirectMessage",
		[]any{"test123", "test test123", "hello", "world"},
		nil)
}

func TestClient_Channels_Tokens_GetAll(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/channels/tokens",
		wantQuery:     map[string]string{"project": "test123"},
		wantResultKey: "tokens",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"GetAll",
		[]any{"test123"},
		[]*types.ChannelToken{{ID: "hello"}})
}
