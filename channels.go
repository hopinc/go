package hop

import (
	"context"
	"net/url"

	"go.hop.io/sdk/types"
	"golang.org/x/sync/errgroup"
)

// Create is used to create a channel. The channelType argument should be the type of channel that you want to create, state
// should be a map that you want to associate with the channel, and id should be the ID you wish to specify (oe a blank string
// if you wish for this to be auto-generated).
func (c ClientCategoryChannels) Create(
	ctx context.Context, channelType types.ChannelType, state map[string]any, id string,
	opts ...ClientOption,
) (*types.Channel, error) {
	// Setup everything we need to do the request.
	method := "POST"
	path := "/channels"
	query := map[string]string{}
	if id != "" {
		method = "PUT"
		path += "/" + url.PathEscape(id)
		query["channel_id"] = id
	}
	if state == nil {
		state = map[string]any{}
	}

	// Do the request.
	var ch types.Channel
	err := c.c.do(ctx, ClientArgs{
		Method:    method,
		Path:      path,
		ResultKey: "channel",
		Query:     query,
		Body:      map[string]any{"type": channelType, "state": state},
		Result:    &ch,
		Ignore404: false,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

// Get is used to get a channel. Will throw types.NotFound if it was not found.
func (c ClientCategoryChannels) Get(ctx context.Context, id string, opts ...ClientOption) (*types.Channel, error) {
	var ch types.Channel
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/channels/" + url.PathEscape(id),
		ResultKey: "channel",
		Result:    &ch,
		Ignore404: false,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

// GetAll returns a paginator to get all the channels.
func (c ClientCategoryChannels) GetAll() *Paginator[*types.Channel] {
	return &Paginator[*types.Channel]{
		c:         c.c,
		total:     -1,
		path:      "/channels",
		resultKey: "channels",
		sortBy:    "created_at",
	}
}

// SubscribeToken is used to subscribe a token to a channel.
func (c ClientCategoryChannels) SubscribeToken(ctx context.Context, channelId, token string, opts ...ClientOption) error {
	path := "/channels/" + url.PathEscape(channelId) + "/subscribers/" + url.PathEscape(token)
	return c.c.do(ctx, ClientArgs{
		Method:    "PUT",
		Path:      path,
		Ignore404: false,
	}, opts)
}

// RemoveToken is used to remove a token from a channel.
func (c ClientCategoryChannels) RemoveToken(ctx context.Context, channelId, token string, opts ...ClientOption) error {
	path := "/channels/" + url.PathEscape(channelId) + "/subscribers/" + url.PathEscape(token)
	return c.c.do(ctx, ClientArgs{
		Method:    "DELETE",
		Path:      path,
		Ignore404: false,
	}, opts)
}

// SubscribeTokens is used to subscribe many tokens to a channel.
func (c ClientCategoryChannels) SubscribeTokens(ctx context.Context, channelId string, tokens []string, opts ...ClientOption) error {
	eg := errgroup.Group{}
	for _, v := range tokens {
		token := v
		eg.Go(func() error {
			return c.SubscribeToken(ctx, channelId, token, opts...)
		})
	}
	return eg.Wait()
}

// SetState is used to set the state of a channel.
func (c ClientCategoryChannels) SetState(ctx context.Context, id string, state map[string]any, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method:    "PUT",
		Path:      "/channels/" + url.PathEscape(id) + "/state",
		Body:      state,
		Ignore404: false,
	}, opts)
}

// PatchState is used to patch the state of a channel.
func (c ClientCategoryChannels) PatchState(ctx context.Context, id string, state map[string]any, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method:    "PATCH",
		Path:      "/channels/" + url.PathEscape(id) + "/state",
		Body:      state,
		Ignore404: false,
	}, opts)
}

// PublishMessage is used to publish an event to the channel.
func (c ClientCategoryChannels) PublishMessage(ctx context.Context, channelId, eventName string, data any, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method:    "POST",
		Path:      "/channels/" + url.PathEscape(channelId) + "/messages",
		Body:      map[string]any{"e": eventName, "d": data},
		Ignore404: false,
	}, opts)
}

// Delete is used to delete a channel.
func (c ClientCategoryChannels) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method:    "DELETE",
		Path:      "/channels/" + url.PathEscape(id),
		Ignore404: false,
	}, opts)
}

// GetStats is used to get the stats of a channel.
func (c ClientCategoryChannels) GetStats(ctx context.Context, id string, opts ...ClientOption) (*types.Stats, error) {
	var s types.Stats
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/channels/" + url.PathEscape(id) + "/stats",
		ResultKey: "stats",
		Result:    &s,
		Ignore404: false,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Delete is used to delete a channel token.
func (t ClientCategoryChannelsTokens) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	return t.c.do(ctx, ClientArgs{
		Method:    "DELETE",
		Path:      "/channels/tokens/" + url.PathEscape(id),
		Ignore404: false,
	}, opts)
}

// Create is used to create a new channel token. State is the map of the state of the token (this can be nil), and
// projectId is the project ID to associate the token with (this can be empty unless it is bearer or PAT auth).
func (t ClientCategoryChannelsTokens) Create(ctx context.Context, state map[string]any, opts ...ClientOption) (*types.ChannelToken, error) {
	if t.c.getProjectId(opts) == "" && t.c.getTokenType() != "ptk" {
		return nil, types.InvalidToken("project ID must be specified when creating a channel token with bearer or PAT auth")
	}
	if state == nil {
		state = map[string]any{}
	}
	var ct types.ChannelToken
	err := t.c.do(ctx, ClientArgs{
		Method:    "POST",
		Path:      "/channels/tokens",
		Body:      map[string]any{"state": state},
		ResultKey: "token",
		Result:    &ct,
		Ignore404: false,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &ct, nil
}

// SetState is used to set the state of a channel token.
func (t ClientCategoryChannelsTokens) SetState(ctx context.Context, id string, state map[string]any, opts ...ClientOption) error {
	return t.c.do(ctx, ClientArgs{
		Method:    "PATCH",
		Path:      "/channels/tokens/" + url.PathEscape(id),
		Body:      map[string]any{"state": state},
		Ignore404: false,
	}, opts)
}

// Get is used to get a token by its ID.
func (t ClientCategoryChannelsTokens) Get(ctx context.Context, id string, opts ...ClientOption) (*types.ChannelToken, error) {
	var ct types.ChannelToken
	err := t.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/channels/tokens/" + url.PathEscape(id),
		ResultKey: "token",
		Result:    &ct,
		Ignore404: false,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &ct, nil
}

// IsOnline is used to check if a token is online.
func (t ClientCategoryChannelsTokens) IsOnline(ctx context.Context, id string, opts ...ClientOption) (bool, error) {
	x, err := t.Get(ctx, id, opts...)
	if err != nil {
		return false, err
	}
	return x.IsOnline, nil
}

// PublishDirectMessage is used to publish an event to the channel token.
func (t ClientCategoryChannelsTokens) PublishDirectMessage(
	ctx context.Context, id, eventName string, data any, opts ...ClientOption,
) error {
	return t.c.do(ctx, ClientArgs{
		Method:    "POST",
		Path:      "/channels/tokens/" + url.PathEscape(id) + "/messages",
		Body:      map[string]any{"e": eventName, "d": data},
		Ignore404: false,
	}, opts)
}

// GetAll gets all the tokens.
func (t ClientCategoryChannelsTokens) GetAll(ctx context.Context, opts ...ClientOption) ([]*types.ChannelToken, error) {
	var a []*types.ChannelToken
	err := t.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/channels/tokens",
		ResultKey: "tokens",
		Result:    &a,
		Ignore404: false,
	}, opts)
	if err != nil {
		return nil, err
	}
	return a, nil
}
