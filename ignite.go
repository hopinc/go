package hopgo

import (
	"context"
	"errors"
	"net/url"

	"github.com/hopinc/hop-go/types"
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
		if c.c.tokenType != "bearer" && c.c.tokenType != "ptk" {
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
