package hopgo

import (
	"context"
	"github.com/jakemakesstuff/pinkypromise/promise"
	"net/url"

	"github.com/hopinc/hop-go/types"
)

// Create is used to create a channel. The channelType argument should be the type of channel that you want to create, state
// should be a map that you want to associate with the channel, id should be the ID you wish to specify (oe a blank string
// if you wish for this to be auto-generated), and projectId should be either a project ID to assign this to or a blank string.
func (c Channels) Create(
	ctx context.Context, channelType types.ChannelType, state map[string]any, id, projectId string,
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
	if projectId != "" {
		query["project"] = projectId
	}

	// Do the request.
	var ch types.Channel
	err := c.c.do(ctx, clientArgs{
		method:    method,
		path:      path,
		resultKey: "channel",
		query:     query,
		body:      map[string]any{"type": channelType, "state": state},
		result:    &ch,
		ignore404: false,
	})
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

// Get is used to get a channel. Will throw types.NotFound if it was not found.
func (c Channels) Get(ctx context.Context, id string) (*types.Channel, error) {
	var ch types.Channel
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/channels/" + url.PathEscape(id),
		resultKey: "channel",
		result:    &ch,
		ignore404: false,
	})
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

// GetAll is used to get all the channels.
func (c Channels) GetAll(ctx context.Context, projectId string) ([]*types.Channel, error) {
	var query map[string]string
	if projectId != "" {
		query = map[string]string{"project": projectId}
	}
	var chs []*types.Channel
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/channels",
		resultKey: "channels",
		result:    &chs,
		query:     query,
		ignore404: false,
	})
	if err != nil {
		return nil, err
	}
	return chs, nil
}

// SubscribeToken is used to subscribe a token to a channel.
func (c Channels) SubscribeToken(ctx context.Context, channelId, token string) error {
	path := "/channels/" + url.PathEscape(channelId) + "/subscribers/" + url.PathEscape(token)
	return c.c.do(ctx, clientArgs{
		method:    "PUT",
		path:      path,
		ignore404: false,
	})
}

// SubscribeTokens is used to subscribe many tokens to a channel.
func (c Channels) SubscribeTokens(ctx context.Context, channelId string, tokens []string) error {
	promises := make([]*promise.Promise[struct{}], len(tokens))
	for i, v := range tokens {
		token := v
		promises[i] = promise.NewFn(func() (struct{}, error) {
			return struct{}{}, c.SubscribeToken(ctx, channelId, token)
		})
	}
	_, err := promise.All(promises...)
	return err
}

// GetAllTokens gets all the tokens associated with a channel.
func (c Channels) GetAllTokens(ctx context.Context, channelId string) ([]string, error) {
	var t []string
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/channels/" + url.PathEscape(channelId) + "/tokens",
		resultKey: "tokens",
		result:    &t,
		ignore404: false,
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

// SetState is used to set the state of a channel.
func (c Channels) SetState(ctx context.Context, id string, state map[string]any) error {
	return c.c.do(ctx, clientArgs{
		method:    "PUT",
		path:      "/channels/" + url.PathEscape(id) + "/state",
		body:      state,
		ignore404: false,
	})
}

// PatchState is used to patch the state of a channel.
func (c Channels) PatchState(ctx context.Context, id string, state map[string]any) error {
	return c.c.do(ctx, clientArgs{
		method:    "PATCH",
		path:      "/channels/" + url.PathEscape(id) + "/state",
		body:      state,
		ignore404: false,
	})
}

// PublishMessage is used to publish an event to the channel.
func (c Channels) PublishMessage(ctx context.Context, channelId, eventName string, data any) error {
	return c.c.do(ctx, clientArgs{
		method:    "POST",
		path:      "/channels/" + url.PathEscape(channelId) + "/messages",
		body:      map[string]any{"e": eventName, "d": data},
		ignore404: false,
	})
}

// Delete is used to delete a channel.
func (c Channels) Delete(ctx context.Context, id string) error {
	return c.c.do(ctx, clientArgs{
		method:    "DELETE",
		path:      "/channels/" + url.PathEscape(id),
		ignore404: false,
	})
}

// GetStats is used to get the stats of a channel.
func (c Channels) GetStats(ctx context.Context, id string) (*types.Stats, error) {
	var s types.Stats
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/channels/" + url.PathEscape(id) + "/stats",
		resultKey: "stats",
		result:    &s,
		ignore404: false,
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Delete is used to delete a channel token.
func (t ChannelsTokens) Delete(ctx context.Context, id string) error {
	return t.c.do(ctx, clientArgs{
		method:    "DELETE",
		path:      "/channels/tokens/" + url.PathEscape(id),
		ignore404: false,
	})
}
