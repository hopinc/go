package hop

import (
	"context"
	"net/url"

	"github.com/hopinc/hop-go/types"
)

// Delete is used to delete a token.
func (c ClientCategoryProjectsTokens) Delete(ctx context.Context, projectId, id string) error {
	if c.c.getTokenType() == "ptk" {
		return types.InvalidToken("project tokens cannot be retrieved with a project token")
	}
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/projects/" + url.PathEscape(projectId) + "/tokens/" + url.PathEscape(id),
	})
}

// GetAll is used to get all tokens associated with a project.
func (c ClientCategoryProjectsTokens) GetAll(ctx context.Context, projectId string) ([]*types.ProjectToken, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("project tokens cannot be retrieved with a project token")
	}
	var tokens []*types.ProjectToken
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/projects/" + url.PathEscape(projectId) + "/tokens",
		resultKey: "project_tokens",
		result:    &tokens,
	})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// Create is used to create a token.
func (c ClientCategoryProjectsTokens) Create(ctx context.Context, projectId string, permissions []types.ProjectPermission) (*types.ProjectToken, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("project tokens cannot be created with a project token")
	}
	if permissions == nil {
		permissions = []types.ProjectPermission{}
	}
	var token types.ProjectToken
	err := c.c.do(ctx, clientArgs{
		method:    "POST",
		path:      "/projects/" + url.PathEscape(projectId) + "/tokens",
		resultKey: "project_token",
		body:      map[string][]types.ProjectPermission{"permissions": permissions},
		result:    &token,
	})
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetAllMembers is used to get all members associated with a project.
func (c ClientCategoryProjects) GetAllMembers(ctx context.Context, projectId string) ([]*types.ProjectMember, error) {
	if c.c.getTokenType() == "ptk" {
		projectId = "@this"
	} else {
		if projectId == "" {
			return nil, types.InvalidToken("project ID must be specified when using bearer authentication to get members")
		}
	}

	var members []*types.ProjectMember
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/projects/" + url.PathEscape(projectId) + "/members",
		resultKey: "members",
		result:    &members,
	})
	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetCurrentMember is used to get the current member associated with a project. You cannot use this method with a project token.
func (c ClientCategoryProjects) GetCurrentMember(ctx context.Context, projectId string) (*types.ProjectMember, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("current member cannot be retrieved with a project token")
	}
	if projectId == "" {
		return nil, types.InvalidToken("project ID must be specified when using bearer authentication to get current member")
	}
	var projectMember *types.ProjectMember
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/projects/" + url.PathEscape(projectId) + "/members/@me",
		resultKey: "project_member",
		query:     getProjectIdParam(projectId),
		result:    &projectMember,
	})
	if err != nil {
		return nil, err
	}
	return projectMember, nil
}

// GetAll is used to get all project secrets.
func (c ClientCategoryProjectsSecrets) GetAll(ctx context.Context, projectId string) ([]*types.ProjectSecret, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("project secrets cannot be retrieved with a project token")
	}
	var secrets []*types.ProjectSecret
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/projects/" + url.PathEscape(projectId) + "/secrets",
		resultKey: "secrets",
		result:    &secrets,
	})
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

// Create is used to create a project secret.
func (c ClientCategoryProjectsSecrets) Create(ctx context.Context, projectId string, key, value string) (*types.ProjectSecret, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("project secrets cannot be created with a project token")
	}
	var secret types.ProjectSecret
	err := c.c.do(ctx, clientArgs{
		method:    "PUT",
		path:      "/projects/" + url.PathEscape(projectId) + "/secrets/" + url.PathEscape(key),
		resultKey: "secret",
		body:      plainText(value),
		result:    &secret,
	})
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// Delete is used to delete a project secret.
func (c ClientCategoryProjectsSecrets) Delete(ctx context.Context, projectId, id string) error {
	if c.c.getTokenType() == "ptk" {
		return types.InvalidToken("project secrets cannot be deleted with a project token")
	}
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/projects/" + url.PathEscape(projectId) + "/secrets/" + url.PathEscape(id),
	})
}
