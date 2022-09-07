package hopgo

import (
	"context"
	"errors"
	"github.com/hopinc/hop-go/types"
	"net/url"
)

// AddDomain is used to add a domain to the gateway. The parameter gatewayId is the ID of the gateway to add the domain to,
// and domain is the full name of the domain.
func (c ClientCategoryIgniteGateways) AddDomain(ctx context.Context, gatewayId string, domain string) error {
	return c.c.do(ctx, clientArgs{
		method: "POST",
		path:   "/ignite/gateways/" + url.PathEscape(gatewayId) + "/domains",
		body:   map[string]any{"domain": domain},
	})
}

// Get is used to get a gateway by its ID.
func (c ClientCategoryIgniteGateways) Get(ctx context.Context, id string) (*types.Gateway, error) {
	var gw types.Gateway
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/gateways/" + url.PathEscape(id),
		resultKey: "gateway",
		result:    &gw,
	})
	if err != nil {
		return nil, err
	}
	return &gw, nil
}

// Create is used to create a deployment.
func (c ClientCategoryIgniteDeployments) Create(ctx context.Context, projectId string, deployment *types.DeploymentConfig) (*types.Deployment, error) {
	if projectId == "" {
		if c.c.tokenType != "ptk" {
			return nil, types.InvalidToken("project ID must be specified when using bearer authentication to make deployments")
		}
	} else {
		if c.c.tokenType != "bearer" && c.c.tokenType != "pat" {
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
		query:     getProjectIdParam(projectId),
		body:      deployment,
		result:    &d,
		ignore404: false,
	})
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Get is used to get a deployment by its ID.
func (c ClientCategoryIgniteDeployments) Get(ctx context.Context, projectId, id string) (*types.Deployment, error) {
	var d types.Deployment
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments/" + url.PathEscape(id),
		query:     getProjectIdParam(projectId),
		resultKey: "deployment",
		result:    &d,
	})
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetByName is used to get a deployment by its name.
func (c ClientCategoryIgniteDeployments) GetByName(ctx context.Context, projectId, name string) (*types.Deployment, error) {
	query := map[string]string{"name": name}
	if projectId != "" {
		query["project"] = projectId
	}

	var d types.Deployment
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments/search",
		query:     query,
		resultKey: "deployment",
		result:    &d,
	})
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetContainers is used to get the containers of a deployment.
func (c ClientCategoryIgniteDeployments) GetContainers(ctx context.Context, projectId, id string) ([]*types.Container, error) {
	var containers []*types.Container
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments/" + url.PathEscape(id) + "/containers",
		query:     getProjectIdParam(projectId),
		resultKey: "containers",
		result:    &containers,
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// GetAll is used to get all deployments.
func (c ClientCategoryIgniteDeployments) GetAll(ctx context.Context, projectId string) ([]*types.Deployment, error) {
	var deployments []*types.Deployment
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments",
		query:     getProjectIdParam(projectId),
		resultKey: "deployments",
		result:    &deployments,
	})
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

// Delete is used to delete a deployment by its ID.
func (c ClientCategoryIgniteDeployments) Delete(ctx context.Context, projectId, id string) error {
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/ignite/deployments/" + url.PathEscape(id),
		query:  getProjectIdParam(projectId),
	})
}

// GetAllGateways is used to get all gateways attached to a deployment.
func (c ClientCategoryIgniteDeployments) GetAllGateways(ctx context.Context, projectId, id string) ([]*types.Gateway, error) {
	var gateways []*types.Gateway
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/ignite/deployments/" + url.PathEscape(id) + "/gateways",
		query:     getProjectIdParam(projectId),
		resultKey: "gateways",
		result:    &gateways,
	})
	if err != nil {
		return nil, err
	}
	return gateways, nil
}

// CreateGateway is used to create a gateway attached to a deployment.
func (c ClientCategoryIgniteDeployments) CreateGateway(ctx context.Context, opts types.GatewayCreationOptions) (*types.Gateway, error) {
	var gw types.Gateway
	err := c.c.do(ctx, clientArgs{
		method:    "POST",
		path:      "/ignite/deployments/" + url.PathEscape(opts.DeploymentID) + "/gateways",
		query:     getProjectIdParam(opts.ProjectID),
		body:      opts,
		resultKey: "gateway",
		result:    &gw,
	})
	if err != nil {
		return nil, err
	}
	return &gw, nil
}

// Delete is used to delete a container by its ID.
func (c ClientCategoryIgniteContainers) Delete(ctx context.Context, projectId, id string) error {
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/ignite/containers/" + url.PathEscape(id),
		query:  getProjectIdParam(projectId),
	})
}

// GetLogs is used to get a paginator for the logs of a container.
func (c ClientCategoryIgniteContainers) GetLogs(projectId, id string, limit int, ascOrder bool) *Paginator[*types.ContainerLog] {
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
		query:       getProjectIdParam(projectId),
	}
}

// Used to update the container state.
func (c ClientCategoryIgniteContainers) updateContainerState(ctx context.Context, projectId, id string, state types.ContainerState) error {
	return c.c.do(ctx, clientArgs{
		method: "PUT",
		path:   "/ignite/containers/" + url.PathEscape(id) + "/state",
		body:   map[string]types.ContainerState{"preferred_state": state},
		query:  getProjectIdParam(projectId),
	})
}

// Stop is used to stop a container by its ID.
func (c ClientCategoryIgniteContainers) Stop(ctx context.Context, projectId, id string) error {
	return c.updateContainerState(ctx, projectId, id, types.ContainerStateStopped)
}

// Start is used to start a container by its ID.
func (c ClientCategoryIgniteContainers) Start(ctx context.Context, projectId, id string) error {
	return c.updateContainerState(ctx, projectId, id, types.ContainerStateRunning)
}

// Create is used to create a container.
func (c ClientCategoryIgniteContainers) Create(ctx context.Context, projectId, deploymentId string) (*types.Container, error) {
	var a []*types.Container
	err := c.c.do(ctx, clientArgs{
		method:    "POST",
		path:      "/ignite/deployments/" + url.PathEscape(deploymentId) + "/containers",
		query:     getProjectIdParam(projectId),
		resultKey: "containers",
		result:    &a,
	})
	if err != nil {
		return nil, err
	}
	return a[0], nil
}
