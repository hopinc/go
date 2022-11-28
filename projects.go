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
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path:   "/projects/" + url.PathEscape(projectId) + "/tokens/" + url.PathEscape(id),
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
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/projects/" + url.PathEscape(projectId) + "/tokens",
		ResultKey: "project_tokens",
		Result:    &tokens,
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
	err := c.c.do(ctx, ClientArgs{
		Method:    "POST",
		Path:      "/projects/" + url.PathEscape(projectId) + "/tokens",
		ResultKey: "project_token",
		Body:      map[string][]types.ProjectPermission{"permissions": permissions},
		Result:    &token,
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
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/projects/" + url.PathEscape(projectId) + "/members",
		ResultKey: "members",
		Result:    &members,
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
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/projects/" + url.PathEscape(projectId) + "/members/@me",
		ResultKey: "project_member",
		Result:    &projectMember,
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
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/projects/" + url.PathEscape(projectId) + "/secrets",
		ResultKey: "secrets",
		Result:    &secrets,
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
	err := c.c.do(ctx, ClientArgs{
		Method:    "PUT",
		Path:      "/projects/" + url.PathEscape(projectId) + "/secrets/" + url.PathEscape(key),
		ResultKey: "secret",
		Body:      PlainText(value),
		Result:    &secret,
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
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path:   "/projects/" + url.PathEscape(projectId) + "/secrets/" + url.PathEscape(id),
	}, opts)
}
