package hop

import (
	"testing"

	"github.com/hopinc/hop-go/types"
	"github.com/stretchr/testify/assert"
)

func TestClient_Ignite_Gateways_AddDomain(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "POST",
		wantPath:      "/ignite/gateways/test%20test/domains",
		wantIgnore404: false,
		wantBody:      map[string]any{"domain": "example.com"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteGateways{c: c},
		"AddDomain",
		[]any{"test test", "example.com"},
		nil)
}

func TestClient_Ignite_Gateways_Get(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/gateways/test%20test",
		wantResultKey: "gateway",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteGateways{c: c},
		"Get",
		[]any{"test test"},
		&types.Gateway{ID: "hello"})
}

func TestClient_Ignite_Deployments_Create(t *testing.T) {
	deploymentConfig := &types.DeploymentConfig{
		Resources: types.Resources{
			RAM: "1gb",
		},
	}
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "POST",
		wantPath:      "/ignite/deployments",
		wantResultKey: "deployment",
		wantIgnore404: false,
		wantQuery:     map[string]string{"project": "test123"},
		wantBody:      deploymentConfig,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"Create",
		[]any{"test123", deploymentConfig},
		&types.Deployment{ID: "hello"})
}

func TestClient_Ignite_Deployments_Get(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/deployments/test%20test",
		wantResultKey: "deployment",
		wantIgnore404: false,
		wantQuery:     map[string]string{"project": "test123"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"Get",
		[]any{"test123", "test test"},
		&types.Deployment{ID: "hello"})
}

func TestClient_Ignite_Deployments_GetByName(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/deployments/search",
		wantResultKey: "deployment",
		wantIgnore404: false,
		wantQuery:     map[string]string{"project": "test123", "name": "test test"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetByName",
		[]any{"test123", "test test"},
		&types.Deployment{ID: "hello"})
}

func TestClient_Ignite_Deployments_GetContainers(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/deployments/test%20test/containers",
		wantResultKey: "containers",
		wantIgnore404: false,
		wantQuery:     map[string]string{"project": "test123"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetContainers",
		[]any{"test123", "test test"},
		[]*types.Container{{ID: "hello"}})
}

func TestClient_Ignite_Deployments_GetAll(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/deployments",
		wantResultKey: "deployments",
		wantIgnore404: false,
		wantQuery:     map[string]string{"project": "test123"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetAll",
		[]any{"test123"},
		[]*types.Deployment{{ID: "hello"}})
}

func TestClient_Ignite_Deployments_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "DELETE",
		wantPath:      "/ignite/deployments/test%20test",
		wantIgnore404: false,
		wantQuery:     map[string]string{"project": "test123"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"Delete",
		[]any{"test123", "test test"},
		nil)
}

func TestClient_Ignite_Deployments_GetAllGateways(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/deployments/test%20test/gateways",
		wantResultKey: "gateways",
		wantIgnore404: false,
		wantQuery:     map[string]string{"project": "test123"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetAllGateways",
		[]any{"test123", "test test"},
		[]*types.Gateway{{ID: "hello"}})
}

func TestClient_Ignite_Deployments_CreateGateway(t *testing.T) {
	gatewayConfig := types.GatewayCreationOptions{
		DeploymentID: "test test",
		Type:         types.GatewayTypeExternal,
	}
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "POST",
		wantPath:      "/ignite/deployments/test%20test/gateways",
		wantResultKey: "gateway",
		wantIgnore404: false,
		wantBody:      gatewayConfig,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"CreateGateway",
		[]any{gatewayConfig},
		&types.Gateway{ID: "hello"})
}

func TestClient_Ignite_Containers_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "DELETE",
		wantPath:      "/ignite/containers/test%20test",
		wantIgnore404: false,
		wantQuery:     map[string]string{"project": "test123"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteContainers{c: c},
		"Delete",
		[]any{"test123", "test test"},
		nil)
}

func TestClient_Ignite_Containers_GetLogs(t *testing.T) {
	c := &mockClientDoer{}
	res := (&ClientCategoryIgniteContainers{c: c}).GetLogs(
		"test123", "test test", 10, true)
	assert.Equal(t, res, &Paginator[*types.ContainerLog]{
		c:           c,
		total:       -1,
		offsetStrat: true,
		limit:       10,
		path:        "/ignite/containers/test%20test/logs",
		resultKey:   "logs",
		sortBy:      "timestamp",
		orderBy:     "asc",
		query:       map[string]string{"project": "test123"},
	})
}

func TestClient_Ignite_Containers_Stop(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "PUT",
		wantPath:      "/ignite/containers/test%20test/state",
		wantIgnore404: false,
		wantBody:      map[string]types.ContainerState{"preferred_state": "stopped"},
		wantQuery:     map[string]string{"project": "test123"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteContainers{c: c},
		"Stop",
		[]any{"test123", "test test"},
		nil)
}

func TestClient_Ignite_Containers_Start(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "PUT",
		wantPath:      "/ignite/containers/test%20test/state",
		wantIgnore404: false,
		wantBody:      map[string]types.ContainerState{"preferred_state": "running"},
		wantQuery:     map[string]string{"project": "test123"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteContainers{c: c},
		"Start",
		[]any{"test123", "test test"},
		nil)
}
