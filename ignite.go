package hop

import (
	"context"
	"errors"
	"net/url"

	"go.hop.io/sdk/types"
)

// AddDomain is used to add a domain to the gateway. The parameter gatewayId is the ID of the gateway to add the domain to,
// and domain is the full name of the domain.
func (c ClientCategoryIgniteGateways) AddDomain(ctx context.Context, gatewayId string, domain string, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method: "POST",
		Path:   "/ignite/gateways/" + url.PathEscape(gatewayId) + "/domains",
		Body:   map[string]any{"domain": domain},
	}, opts)
}

// GetDomain is used to get a domain by its ID.
func (c ClientCategoryIgniteGateways) GetDomain(
	ctx context.Context, domainId string, opts ...ClientOption,
) (*types.Domain, error) {
	var d types.Domain
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/ignite/domains/" + url.PathEscape(domainId),
		ResultKey: "domain",
		Result:    &d,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// DeleteDomain is used to delete a domain by its ID.
func (c ClientCategoryIgniteGateways) DeleteDomain(ctx context.Context, domainId string, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path:   "/ignite/domains/" + url.PathEscape(domainId),
	}, opts)
}

// Get is used to get a gateway by its ID.
func (c ClientCategoryIgniteGateways) Get(ctx context.Context, id string, opts ...ClientOption) (*types.Gateway, error) {
	var gw types.Gateway
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/ignite/gateways/" + url.PathEscape(id),
		ResultKey: "gateway",
		Result:    &gw,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &gw, nil
}

// Delete is used to delete a gateway by its ID.
func (c ClientCategoryIgniteGateways) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path:   "/ignite/gateways/" + url.PathEscape(id),
	}, opts)
}

// Update is used to update a gateway by its ID.
func (c ClientCategoryIgniteGateways) Update(
	ctx context.Context, id string, updateOpts types.IgniteGatewayUpdateOpts, opts ...ClientOption,
) (*types.Gateway, error) {
	var gw types.Gateway
	err := c.c.do(ctx, ClientArgs{
		Method:    "PATCH",
		Path:      "/ignite/gateways/" + url.PathEscape(id),
		ResultKey: "gateway",
		Body:      updateOpts,
		Result:    &gw,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &gw, nil
}

// Create is used to create a deployment.
func (c ClientCategoryIgniteDeployments) Create(
	ctx context.Context, deployment *types.DeploymentConfig, opts ...ClientOption,
) (*types.Deployment, error) {
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
	err = c.c.do(ctx, ClientArgs{
		Method:    "POST",
		Path:      "/ignite/deployments",
		ResultKey: "deployment",
		Body:      deployment,
		Result:    &d,
		Ignore404: false,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Get is used to get a deployment by its ID.
func (c ClientCategoryIgniteDeployments) Get(ctx context.Context, id string, opts ...ClientOption) (*types.Deployment, error) {
	var d types.Deployment
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/ignite/deployments/" + url.PathEscape(id),
		ResultKey: "deployment",
		Result:    &d,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetByName is used to get a deployment by its name.
func (c ClientCategoryIgniteDeployments) GetByName(ctx context.Context, name string, opts ...ClientOption) (*types.Deployment, error) {
	var d types.Deployment
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/ignite/deployments/search",
		Query:     map[string]string{"name": name},
		ResultKey: "deployment",
		Result:    &d,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Update is used to update a deployment by its ID.
func (c ClientCategoryIgniteDeployments) Update(
	ctx context.Context, id string, updateOpts types.IgniteDeploymentUpdateOpts, opts ...ClientOption,
) (*types.Deployment, error) {
	var d types.Deployment
	err := c.c.do(ctx, ClientArgs{
		Method:    "PATCH",
		Path:      "/ignite/deployments/" + url.PathEscape(id),
		ResultKey: "deployment",
		Body:      updateOpts,
		Result:    &d,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Patch is used to patch a deployment by its ID.
//
// Deprecated: use Update instead.
func (c ClientCategoryIgniteDeployments) Patch(
	ctx context.Context, id string, patchOpts types.IgniteDeploymentPatchOpts, opts ...ClientOption,
) (*types.Deployment, error) {
	return c.Update(ctx, id, patchOpts, opts...)
}

// GetContainers is used to get the containers of a deployment.
func (c ClientCategoryIgniteDeployments) GetContainers(ctx context.Context, id string, opts ...ClientOption) ([]*types.Container, error) {
	var containers []*types.Container
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/ignite/deployments/" + url.PathEscape(id) + "/containers",
		ResultKey: "containers",
		Result:    &containers,
	}, opts)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// GetAll is used to get all deployments.
func (c ClientCategoryIgniteDeployments) GetAll(ctx context.Context, opts ...ClientOption) ([]*types.Deployment, error) {
	var deployments []*types.Deployment
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/ignite/deployments",
		ResultKey: "deployments",
		Result:    &deployments,
	}, opts)
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

// Delete is used to delete a deployment by its ID.
func (c ClientCategoryIgniteDeployments) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path:   "/ignite/deployments/" + url.PathEscape(id),
	}, opts)
}

// GetAllGateways is used to get all gateways attached to a deployment.
func (c ClientCategoryIgniteDeployments) GetAllGateways(ctx context.Context, id string, opts ...ClientOption) ([]*types.Gateway, error) {
	var gateways []*types.Gateway
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/ignite/deployments/" + url.PathEscape(id) + "/gateways",
		ResultKey: "gateways",
		Result:    &gateways,
	}, opts)
	if err != nil {
		return nil, err
	}
	return gateways, nil
}

// CreateGateway is used to create a gateway attached to a deployment.
func (c ClientCategoryIgniteDeployments) CreateGateway(
	ctx context.Context, opts types.GatewayCreationOptions, clientOpts ...ClientOption,
) (*types.Gateway, error) {
	if opts.ProjectID != "" { //nolint:staticcheck // Support deprecated field.
		clientOpts = append(clientOpts, WithProjectID(opts.ProjectID)) //nolint:staticcheck // Support deprecated field.
	}

	var gw types.Gateway
	err := c.c.do(ctx, ClientArgs{
		Method:    "POST",
		Path:      "/ignite/deployments/" + url.PathEscape(opts.DeploymentID) + "/gateways",
		Body:      opts,
		ResultKey: "gateway",
		Result:    &gw,
	}, clientOpts)
	if err != nil {
		return nil, err
	}
	return &gw, nil
}

// Delete is used to delete a container by its ID.
func (c ClientCategoryIgniteContainers) Delete(ctx context.Context, id string, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path:   "/ignite/containers/" + url.PathEscape(id),
	}, opts)
}

// DeleteAndRecreate is used to delete a container by its ID and then recreate it. This is another function to avoid a
// breaking change.
func (c ClientCategoryIgniteContainers) DeleteAndRecreate(ctx context.Context, id string, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path:   "/ignite/containers/" + url.PathEscape(id),
		Query:  map[string]string{"recreate": "true"},
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
func (c ClientCategoryIgniteContainers) updateContainerState(
	ctx context.Context, id string, state types.ContainerState, opts []ClientOption,
) error {
	return c.c.do(ctx, ClientArgs{
		Method: "PUT",
		Path:   "/ignite/containers/" + url.PathEscape(id) + "/state",
		Body:   map[string]types.ContainerState{"preferred_state": state},
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
	err := c.c.do(ctx, ClientArgs{
		Method:    "POST",
		Path:      "/ignite/deployments/" + url.PathEscape(deploymentId) + "/containers",
		ResultKey: "containers",
		Result:    &a,
	}, opts)
	if err != nil {
		return nil, err
	}
	return a[0], nil
}

// Scale is used to scale the container count of a deployment.
func (c ClientCategoryIgniteDeployments) Scale(
	ctx context.Context, deploymentId string, containerCount uint, opts ...ClientOption,
) ([]*types.Container, error) {
	var a []*types.Container
	err := c.c.do(ctx, ClientArgs{
		Method:    "PATCH",
		Path:      "/ignite/deployments/" + url.PathEscape(deploymentId) + "/scale",
		Body:      map[string]uint{"scale": containerCount},
		ResultKey: "containers",
		Result:    &a,
	}, opts)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// NewHealthCheck is used to set a health check on a deployment. Returns the health check ID.
func (c ClientCategoryIgniteDeployments) NewHealthCheck(
	ctx context.Context, createOpts types.HealthCheckCreateOpts, opts ...ClientOption,
) (*types.HealthCheck, error) {
	// Get the deployment ID.
	deploymentId := createOpts.DeploymentID
	createOpts.DeploymentID = ""

	// Set the defaults.
	if createOpts.Protocol == "" {
		createOpts.Protocol = types.HealthCheckProtocolHTTP
	}
	if createOpts.Path == "" {
		createOpts.Path = "/"
	}
	if createOpts.Port == 0 {
		createOpts.Port = 8080
	}
	if createOpts.InitialDelay == 0 {
		createOpts.InitialDelay = types.SecondsFromInt(5)
	}
	if createOpts.Interval == 0 {
		createOpts.Interval = types.SecondsFromInt(60)
	}
	if createOpts.Timeout == 0 {
		createOpts.Timeout = types.MillisecondsFromInt(50)
	}
	if createOpts.MaxRetries == 0 {
		createOpts.MaxRetries = 3
	}

	// Do the HTTP request.
	var res types.HealthCheck
	err := c.c.do(ctx, ClientArgs{
		Method:    "POST",
		Path:      "/ignite/deployments/" + url.PathEscape(deploymentId) + "/health-checks",
		Body:      createOpts,
		ResultKey: "health_check",
		Result:    &res,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetHealthChecks is used to get the health checks attached to a deployment ID.
func (c ClientCategoryIgniteDeployments) GetHealthChecks(
	ctx context.Context, deploymentId string, opts ...ClientOption,
) ([]*types.HealthCheck, error) {
	var res []*types.HealthCheck
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/ignite/deployments/" + url.PathEscape(deploymentId) + "/health-checks",
		ResultKey: "health_checks",
		Result:    &res,
	}, opts)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteHealthCheck is used to delete a health check by its ID.
func (c ClientCategoryIgniteDeployments) DeleteHealthCheck(
	ctx context.Context, deploymentId, healthCheckId string, opts ...ClientOption,
) error {
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path: "/ignite/deployments/" + url.PathEscape(deploymentId) + "/health-checks/" +
			url.PathEscape(healthCheckId),
	}, opts)
}

// UpdateHealthCheck is used to update a health check.
func (c ClientCategoryIgniteDeployments) UpdateHealthCheck(
	ctx context.Context, updateOpts types.HealthCheckUpdateOpts, opts ...ClientOption,
) (*types.HealthCheck, error) {
	var res types.HealthCheck
	err := c.c.do(ctx, ClientArgs{
		Method: "PATCH",
		Path: "/ignite/deployments/" + url.PathEscape(updateOpts.DeploymentID) + "/health-checks/" +
			url.PathEscape(updateOpts.HealthCheckID),
		Body:      updateOpts,
		ResultKey: "health_check",
		Result:    &res,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// HealthCheckStates is used to get the state of health checks for a deployment.
func (c ClientCategoryIgniteDeployments) HealthCheckStates(
	ctx context.Context, deploymentId string, opts ...ClientOption,
) ([]*types.HealthCheckState, error) {
	var res []*types.HealthCheckState
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/ignite/deployments/" + url.PathEscape(deploymentId) + "/health-check-state",
		ResultKey: "health_check_states",
		Result:    &res,
	}, opts)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetStorageStats is used to get stats about storage.
func (c ClientCategoryIgniteDeployments) GetStorageStats(
	ctx context.Context, deploymentId string, opts ...ClientOption,
) (types.DeploymentStorageInfo, error) {
	var res types.DeploymentStorageInfo
	err := c.c.do(ctx, ClientArgs{
		Method: "GET",
		Path:   "/ignite/deployments/" + url.PathEscape(deploymentId) + "/storage",
		Result: &res,
	}, opts)
	if err != nil {
		return types.DeploymentStorageInfo{}, err
	}
	return res, nil
}
