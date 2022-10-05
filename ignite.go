package hop

import (
	"context"
	"errors"
	"go.hop.io/sdk/types"
	"net/url"
)

// AddDomain is used to add a domain to the gateway. The parameter gatewayId is the ID of the gateway to add the domain to,
// and domain is the full name of the domain.
func (c ClientCategoryIgniteGateways) AddDomain(ctx context.Context, gatewayId string, domain string, opts ...ClientOption) error {
	return c.c.do(ctx, clientArgs{
		method: "POST",
		path:   "/ignite/gateways/" + url.PathEscape(gatewayId) + "/domains",
		body:   map[string]any{"domain": domain},
	}, opts)
}

// Get is used to get a gateway by its ID.
func (c ClientCategoryIgniteGateways) Get(ctx context.Context, id string, opts ...ClientOption) (*types.Gateway, error) {
	var gw types.Gateway
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/gateways/" + url.PathEscape(id),
		resultKey: "gateway",
		result:    &gw,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &gw, nil
}

// Create is used to create a deployment.
func (c ClientCategoryIgniteDeployments) Create(ctx context.Context, deployment *types.DeploymentConfig, opts ...ClientOption) (*types.Deployment, error) {
	if c.c.getProjectId(opts) == "" {
		if c.c.getTokenType() != "ptk" {
			return nil, types.InvalidToken("project ID must be specified when using bearer authentication to make deployments")
		}
	} else {
		if c.c.getTokenType() != "bearer" && c.c.getTokenType() != "pat" {
			return nil, types.InvalidToken("project ID must not be specified if it is implied")
		}
	}

	ramSize, err := deployment.Resources.RAM.Bytes()
	if err != nil {
		return nil, err
	}
	if 6e+6 >= ramSize {
		return nil, errors.New("ram must be at least 6MB")
	}

	var d types.Deployment
	err = c.c.do(ctx, clientArgs{
		method:    "POST",
		path:      "/ignite/deployments",
		resultKey: "deployment",
		body:      deployment,
		result:    &d,
		ignore404: false,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Get is used to get a deployment by its ID.
func (c ClientCategoryIgniteDeployments) Get(ctx context.Context, id string, opts ...ClientOption) (*types.Deployment, error) {
	var d types.Deployment
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments/" + url.PathEscape(id),
		resultKey: "deployment",
		result:    &d,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetByName is used to get a deployment by its name.
func (c ClientCategoryIgniteDeployments) GetByName(ctx context.Context, name string, opts ...ClientOption) (*types.Deployment, error) {
	var d types.Deployment
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments/search",
		query:     map[string]string{"name": name},
		resultKey: "deployment",
		result:    &d,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Patch is used to patch a deployment by its ID.
func (c ClientCategoryIgniteDeployments) Patch(ctx context.Context, id string, patchOpts types.IgniteDeploymentPatchOpts, opts ...ClientOption) (*types.Deployment, error) {
	var d types.Deployment
	err := c.c.do(ctx, clientArgs{
		method:    "PATCH",
		path:      "/ignite/deployments/" + url.PathEscape(id),
		resultKey: "deployment",
		body:      patchOpts,
		result:    &d,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetContainers is used to get the containers of a deployment.
func (c ClientCategoryIgniteDeployments) GetContainers(ctx context.Context, id string, opts ...ClientOption) ([]*types.Container, error) {
	var containers []*types.Container
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments/" + url.PathEscape(id) + "/containers",
		resultKey: "containers",
		result:    &containers,
	}, opts)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// GetAll is used to get all deployments.
func (c ClientCategoryIgniteDeployments) GetAll(ctx context.Context, opts ...ClientOption) ([]*types.Deployment, error) {
	var deployments []*types.Deployment
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments",
		resultKey: "deployments",
		result:    &deployments,
	}, opts)
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

// Delete is used to delete a deployment by its ID.
func (c ClientCategoryIgniteDeployments) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/ignite/deployments/" + url.PathEscape(id),
	}, opts)
}

// GetAllGateways is used to get all gateways attached to a deployment.
func (c ClientCategoryIgniteDeployments) GetAllGateways(ctx context.Context, id string, opts ...ClientOption) ([]*types.Gateway, error) {
	var gateways []*types.Gateway
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments/" + url.PathEscape(id) + "/gateways",
		resultKey: "gateways",
		result:    &gateways,
	}, opts)
	if err != nil {
		return nil, err
	}
	return gateways, nil
}

// CreateGateway is used to create a gateway attached to a deployment.
func (c ClientCategoryIgniteDeployments) CreateGateway(ctx context.Context, opts types.GatewayCreationOptions, clientOpts ...ClientOption) (*types.Gateway, error) {
	var gw types.Gateway
	err := c.c.do(ctx, clientArgs{
		method:    "POST",
		path:      "/ignite/deployments/" + url.PathEscape(opts.DeploymentID) + "/gateways",
		body:      opts,
		resultKey: "gateway",
		result:    &gw,
	}, clientOpts)
	if err != nil {
		return nil, err
	}
	return &gw, nil
}

// Delete is used to delete a container by its ID.
func (c ClientCategoryIgniteContainers) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/ignite/containers/" + url.PathEscape(id),
	}, opts)
}

// GetLogs is used to get a paginator for the logs of a container.
func (c ClientCategoryIgniteContainers) GetLogs(id string, limit int, ascOrder bool) *Paginator[*types.ContainerLog] {
	orderBy := "desc"
	if ascOrder {
		orderBy = "asc"
	}
	return &Paginator[*types.ContainerLog]{
		c:           c.c,
		total:       -1,
		offsetStrat: true,
		limit:       limit,
		path:        "/ignite/containers/" + url.PathEscape(id) + "/logs",
		resultKey:   "logs",
		sortBy:      "timestamp",
		orderBy:     orderBy,
	}
}

// Used to update the container state.
func (c ClientCategoryIgniteContainers) updateContainerState(ctx context.Context, id string, state types.ContainerState, opts []ClientOption) error {
	return c.c.do(ctx, clientArgs{
		method: "PUT",
		path:   "/ignite/containers/" + url.PathEscape(id) + "/state",
		body:   map[string]types.ContainerState{"preferred_state": state},
	}, opts)
}

// Stop is used to stop a container by its ID.
func (c ClientCategoryIgniteContainers) Stop(ctx context.Context, id string, opts ...ClientOption) error {
	return c.updateContainerState(ctx, id, types.ContainerStateStopped, opts)
}

// Start is used to start a container by its ID.
func (c ClientCategoryIgniteContainers) Start(ctx context.Context, id string, opts ...ClientOption) error {
	return c.updateContainerState(ctx, id, types.ContainerStateRunning, opts)
}

// Create is used to create a container.
func (c ClientCategoryIgniteContainers) Create(ctx context.Context, deploymentId string, opts ...ClientOption) (*types.Container, error) {
	var a []*types.Container
	err := c.c.do(ctx, clientArgs{
		method:    "POST",
		path:      "/ignite/deployments/" + url.PathEscape(deploymentId) + "/containers",
		resultKey: "containers",
		result:    &a,
	}, opts)
	if err != nil {
		return nil, err
	}
	return a[0], nil
}
