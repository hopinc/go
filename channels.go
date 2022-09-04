package hopgo

import (
	"context"
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
