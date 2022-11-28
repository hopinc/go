package hop

import (
	"context"
	"net/url"

	"go.hop.io/sdk/types"
)

// GetAll is used to get all rooms associated with a pipe.
func (c ClientCategoryPipeRooms) GetAll(ctx context.Context, opts ...ClientOption) ([]*types.Room, error) {
	if c.c.getProjectId(opts) == "" && c.c.getTokenType() != "ptk" {
		return nil, types.InvalidToken("project ID must be specified when using bearer authentication to get rooms")
	}
	var rooms []*types.Room
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/pipe/rooms",
		ResultKey: "rooms",
		Result:    &rooms,
	}, opts)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// Create is used to create a room.
func (c ClientCategoryPipeRooms) Create(
	ctx context.Context, opts types.RoomCreationOptions, clientOpts ...ClientOption,
) (*types.Room, error) {
	var room types.Room
	err := c.c.do(ctx, ClientArgs{
		Method:    "POST",
		Path:      "/pipe/rooms",
		ResultKey: "room",
		Result:    &room,
		Body:      opts,
	}, clientOpts)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// Delete is used to delete a room.
func (c ClientCategoryPipeRooms) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path:   "/pipe/rooms/" + url.PathEscape(id),
	}, opts)
}
