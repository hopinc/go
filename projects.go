package hop

import (
	"context"
	"net/url"

	"go.hop.io/sdk/types"
)

// Delete is used to delete a token. The project ID MUST be specified in client options (either at a client or function level).
func (c ClientCategoryProjectsTokens) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	if c.c.getTokenType() == "ptk" {
		return types.InvalidToken("project tokens cannot be retrieved with a project token")
	}
	projectId := c.c.getProjectId(opts)
	if projectId == "" {
		return types.InvalidToken("project id must be specified")
	}
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/projects/" + url.PathEscape(projectId) + "/tokens/" + url.PathEscape(id),
	}, opts)
}

// GetAll is used to get all tokens associated with a project. The project ID MUST be specified in client options
// (either at a client or function level).
func (c ClientCategoryProjectsTokens) GetAll(ctx context.Context, opts ...ClientOption) ([]*types.ProjectToken, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("project tokens cannot be retrieved with a project token")
	}
	projectId := c.c.getProjectId(opts)
	if projectId == "" {
		return nil, types.InvalidToken("project id must be specified")
	}
	var tokens []*types.ProjectToken
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/projects/" + url.PathEscape(projectId) + "/tokens",
		resultKey: "project_tokens",
		result:    &tokens,
	}, opts)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// Create is used to create a token. The project ID MUST be specified in client options (either at a client or function level).
func (c ClientCategoryProjectsTokens) Create(
	ctx context.Context, permissions []types.ProjectPermission, opts ...ClientOption,
) (*types.ProjectToken, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("project tokens cannot be created with a project token")
	}
	projectId := c.c.getProjectId(opts)
	if projectId == "" {
		return nil, types.InvalidToken("project id must be specified")
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
	}, opts)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetAllMembers is used to get all members associated with a project. The project ID MUST be specified in client
// options (either at a client or function level).
func (c ClientCategoryProjects) GetAllMembers(ctx context.Context, opts ...ClientOption) ([]*types.ProjectMember, error) {
	projectId := c.c.getProjectId(opts)
	if c.c.getTokenType() == "ptk" {
		projectId = "@this"
	} else if projectId == "" {
		return nil, types.InvalidToken("project ID must be specified when using bearer authentication to get members")
	}

	var members []*types.ProjectMember
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/projects/" + url.PathEscape(projectId) + "/members",
		resultKey: "members",
		result:    &members,
	}, opts)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetCurrentMember is used to get the current member associated with a project. You cannot use this method with a project token.
// The project ID MUST be specified in client options (either at a client or function level).
func (c ClientCategoryProjects) GetCurrentMember(ctx context.Context, opts ...ClientOption) (*types.ProjectMember, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("current member cannot be retrieved with a project token")
	}
	projectId := c.c.getProjectId(opts)
	if projectId == "" {
		return nil, types.InvalidToken("project ID must be specified when using bearer authentication to get current member")
	}
	var projectMember types.ProjectMember
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/projects/" + url.PathEscape(projectId) + "/members/@me",
		resultKey: "project_member",
		result:    &projectMember,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &projectMember, nil
}

// GetAll is used to get all project secrets. The project ID MUST be specified in client options (either at a client or function level).
func (c ClientCategoryProjectsSecrets) GetAll(ctx context.Context, opts ...ClientOption) ([]*types.ProjectSecret, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("project secrets cannot be retrieved with a project token")
	}
	projectId := c.c.getProjectId(opts)
	if projectId == "" {
		return nil, types.InvalidToken("project ID must be specified when using bearer authentication to get all secrets")
	}
	var secrets []*types.ProjectSecret
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/projects/" + url.PathEscape(projectId) + "/secrets",
		resultKey: "secrets",
		result:    &secrets,
	}, opts)
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

// Create is used to create a project secret. The project ID MUST be specified in client options (either at a client or function level).
func (c ClientCategoryProjectsSecrets) Create(ctx context.Context, key, value string, opts ...ClientOption) (*types.ProjectSecret, error) {
	if c.c.getTokenType() == "ptk" {
		return nil, types.InvalidToken("project secrets cannot be created with a project token")
	}
	projectId := c.c.getProjectId(opts)
	if projectId == "" {
		return nil, types.InvalidToken("project ID must be specified when using bearer authentication to get all secrets")
	}
	var secret types.ProjectSecret
	err := c.c.do(ctx, clientArgs{
		method:    "PUT",
		path:      "/projects/" + url.PathEscape(projectId) + "/secrets/" + url.PathEscape(key),
		resultKey: "secret",
		body:      plainText(value),
		result:    &secret,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// Delete is used to delete a project secret.
func (c ClientCategoryProjectsSecrets) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	if c.c.getTokenType() == "ptk" {
		return types.InvalidToken("project secrets cannot be deleted with a project token")
	}
	projectId := c.c.getProjectId(opts)
	if projectId == "" {
		return types.InvalidToken("project ID must be specified when using bearer authentication to get all secrets")
	}
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/projects/" + url.PathEscape(projectId) + "/secrets/" + url.PathEscape(id),
	}, opts)
}
