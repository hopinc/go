package hopgo

import (
	"context"
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
