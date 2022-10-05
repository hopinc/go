package hop

import (
	"testing"

	"go.hop.io/sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestClient_Channels_Create(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "POST",
		wantPath:       "/channels",
		wantResultKey:  "channel",
		wantIgnore404:  false,
		wantQuery:      map[string]string{},
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantBody:       map[string]any{"type": types.ChannelTypePrivate, "state": map[string]any{"do_you_like_waffles": true}},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"Create",
		[]any{types.ChannelTypePrivate, map[string]any{"do_you_like_waffles": true}, "", WithProjectID("test123")},
		&types.Channel{ChannelPartial: types.ChannelPartial{ID: "hello"}})
}

func TestClient_Channels_Get(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/channels/test%20test",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantResultKey:  "channel",
		wantIgnore404:  false,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"Get",
		[]any{"test test", WithProjectID("test123")},
		&types.Channel{ChannelPartial: types.ChannelPartial{ID: "hello"}})
}

func TestClient_Channels_GetAll(t *testing.T) {
	c := &mockClientDoer{}
	res := (&ClientCategoryChannels{c: c}).GetAll()
	assert.Equal(t, res, &Paginator[*types.Channel]{
		c:         c,
		total:     -1,
		path:      "/channels",
		resultKey: "channels",
		sortBy:    "created_at",
	})
}

func TestClient_Channels_SubscribeToken(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "PUT",
		wantPath:       "/channels/test%20test/subscribers/testing%20testing123",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantIgnore404:  false,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"SubscribeToken",
		[]any{"test test", "testing testing123", WithProjectID("test123")},
		nil)
}

func TestClient_Channels_SetState(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "PUT",
		wantPath:       "/channels/test%20test/state",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantIgnore404:  false,
		wantBody:       map[string]any{"do_you_like_waffles": true},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"SetState",
		[]any{"test test", map[string]any{"do_you_like_waffles": true}, WithProjectID("test123")},
		nil)
}

func TestClient_Channels_PatchState(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "PATCH",
		wantPath:       "/channels/test%20test/state",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantIgnore404:  false,
		wantBody:       map[string]any{"do_you_like_waffles": true},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"PatchState",
		[]any{"test test", map[string]any{"do_you_like_waffles": true}, WithProjectID("test123")},
		nil)
}

func TestClient_Channels_PublishMessage(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "POST",
		wantPath:       "/channels/test%20test/messages",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantIgnore404:  false,
		wantBody:       map[string]any{"e": "hello", "d": "world"},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"PublishMessage",
		[]any{"test test", "hello", "world", WithProjectID("test123")},
		nil)
}

func TestClient_Channels_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "DELETE",
		wantPath:       "/channels/test%20test",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantIgnore404:  false,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"Delete",
		[]any{"test test", WithProjectID("test123")},
		nil)
}

func TestClient_Channels_GetStats(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/channels/test%20test/stats",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantResultKey:  "stats",
		wantIgnore404:  false,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannels{c: c},
		"GetStats",
		[]any{"test test", WithProjectID("test123")},
		&types.Stats{OnlineCount: 9000})
}

func TestClient_Channels_Tokens_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "DELETE",
		wantPath:       "/channels/tokens/test%20test123",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantIgnore404:  false,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"Delete",
		[]any{"test test123", WithProjectID("test123")},
		nil)
}

func TestClient_Channels_Tokens_Create(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "POST",
		wantPath:       "/channels/tokens",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantResultKey:  "token",
		wantIgnore404:  false,
		wantBody:       map[string]any{"state": map[string]any{"hello": "world"}},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"Create",
		[]any{map[string]any{"hello": "world"}, WithProjectID("test123")},
		&types.ChannelToken{ID: "hello"})
}

func TestClient_Channels_Tokens_SetState(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "PATCH",
		wantPath:       "/channels/tokens/test%20test123",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantIgnore404:  false,
		wantBody:       map[string]any{"state": map[string]any{"do_you_like_waffles": true}},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"SetState",
		[]any{"test test123", map[string]any{"do_you_like_waffles": true}, WithProjectID("test123")},
		nil)
}

func TestClient_Channels_Tokens_Get(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/channels/tokens/test%20test123",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantResultKey:  "token",
		wantIgnore404:  false,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"Get",
		[]any{"test test123", WithProjectID("test123")},
		&types.ChannelToken{ID: "hello"})
}

func TestClient_Channels_Tokens_PublishDirectMessage(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "POST",
		wantPath:       "/channels/tokens/test%20test123/messages",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantIgnore404:  false,
		wantBody:       map[string]any{"e": "hello", "d": "world"},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"PublishDirectMessage",
		[]any{"test test123", "hello", "world", WithProjectID("test123")},
		nil)
}

func TestClient_Channels_Tokens_GetAll(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/channels/tokens",
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantResultKey:  "tokens",
		wantIgnore404:  false,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryChannelsTokens{c: c},
		"GetAll",
		[]any{WithProjectID("test123")},
		[]*types.ChannelToken{{ID: "hello"}})
}
