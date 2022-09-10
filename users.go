package hop

import (
	"context"

	"github.com/hopinc/hop-go/types"
)

// Get is used to get the current user.
func (c ClientCategoryUsersMe) Get(ctx context.Context) (*types.UserMeInfo, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("cannot get user with project token")
	}
	var u types.UserMeInfo
	err := c.c.do(ctx, clientArgs{
		method: "GET",
		path:   "/users/@me",
		result: &u,
	})
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreatePat is used to create a personal access token for the current user.
func (c ClientCategoryUsersMe) CreatePat(ctx context.Context, name string) (*types.UserPat, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("cannot create users tokens with project token")
	}
	var pat types.UserPat
	err := c.c.do(ctx, clientArgs{
		method:    "POST",
		path:      "/users/@me/pats",
		resultKey: "pat",
		body:      map[string]string{"name": name},
		result:    &pat,
	})
	if err != nil {
		return nil, err
	}
	return &pat, nil
}

// GetAllPats is used to get all personal access tokens for the current user.
func (c ClientCategoryUsersMe) GetAllPats(ctx context.Context) ([]*types.UserPat, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("cannot get users tokens with project token")
	}
	var pats []*types.UserPat
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/users/@me/pats",
		resultKey: "pats",
		result:    &pats,
	})
	if err != nil {
		return nil, err
	}
	return pats, nil
}

// DeletePat is used to delete a personal access token for the current user.
func (c ClientCategoryUsersMe) DeletePat(ctx context.Context, id string) error {
	if c.c.getTokenType() == "ptk" {
		return types.InvalidToken("cannot delete users tokens with project token")
	}
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/users/@me/pats/" + id,
	})
}
