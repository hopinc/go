package hopgo

import (
	"context"
	"net/url"

	"github.com/hopinc/hop-go/types"
)

// GetAll is used to get all rooms associated with a pipe.
func (c ClientCategoryPipeRooms) GetAll(ctx context.Context, projectId string) ([]*types.Room, error) {
	if projectId == "" && c.c.tokenType != "ptk" {
		return nil, types.InvalidToken("project ID must be specified when using bearer authentication to get rooms")
	}
	var rooms []*types.Room
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/pipe/rooms",
		resultKey: "rooms",
		result:    &rooms,
		query:     getProjectIdParam(projectId),
	})
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// Create is used to create a room.
func (c ClientCategoryPipeRooms) Create(ctx context.Context, opts types.RoomCreationOptions) (*types.Room, error) {
	var room types.Room
	err := c.c.do(ctx, clientArgs{
		method:    "POST",
		path:      "/pipe/rooms",
		resultKey: "room",
		result:    &room,
		body:      opts,
	})
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// Delete is used to delete a room.
func (c ClientCategoryPipeRooms) Delete(ctx context.Context, projectId, id string) error {
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/pipe/rooms/" + url.PathEscape(id),
		query:  getProjectIdParam(projectId),
	})
}
