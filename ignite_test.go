package hop

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.hop.io/sdk/types"
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

func TestClient_Ignite_Gateways_GetDomain(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/domains/test%20test",
		wantResultKey: "domain",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteGateways{c: c},
		"GetDomain",
		[]any{"test test"},
		&types.Domain{Domain: "google.com"})
}

func TestClient_Ignite_Gateways_DeleteDomain(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "DELETE",
		wantPath:      "/ignite/domains/test%20test",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteGateways{c: c},
		"DeleteDomain",
		[]any{"test test"},
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

func TestClient_Ignite_Gateways_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "DELETE",
		wantPath:      "/ignite/gateways/test%20test",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteGateways{c: c},
		"Delete",
		[]any{"test test"},
		nil)
}

func TestClient_Ignite_Gateways_Update(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "PATCH",
		wantPath:      "/ignite/gateways/test%20test",
		wantResultKey: "gateway",
		wantIgnore404: false,
		wantBody:      types.IgniteGatewayUpdateOpts{Name: "new name"},
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteGateways{c: c},
		"Update",
		[]any{"test test", types.IgniteGatewayUpdateOpts{Name: "new name"}},
		&types.Gateway{ID: "hello"})
}

func TestClient_Ignite_Deployments_Create(t *testing.T) {
	deploymentConfig := &types.DeploymentConfig{
		DeploymentConfigPartial: types.DeploymentConfigPartial{
			Resources: types.Resources{
				RAM: "1gb",
			},
		},
	}
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "POST",
		wantPath:       "/ignite/deployments",
		wantResultKey:  "deployment",
		wantIgnore404:  false,
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantBody:       deploymentConfig,
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"Create",
		[]any{deploymentConfig, WithProjectID("test123")},
		&types.Deployment{ID: "hello"})
}

func TestClient_Ignite_Deployments_Get(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/ignite/deployments/test%20test",
		wantResultKey:  "deployment",
		wantIgnore404:  false,
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"Get",
		[]any{"test test", WithProjectID("test123")},
		&types.Deployment{ID: "hello"})
}

func TestClient_Ignite_Deployments_GetByName(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/ignite/deployments/search",
		wantResultKey:  "deployment",
		wantIgnore404:  false,
		wantQuery:      map[string]string{"name": "test test"},
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetByName",
		[]any{"test test", WithProjectID("test123")},
		&types.Deployment{ID: "hello"})
}

func testIgniteUpdate(t *testing.T, name string) {
	t.Helper()
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "PATCH",
		wantPath:       "/ignite/deployments/test%20test",
		wantResultKey:  "deployment",
		wantIgnore404:  false,
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantBody:       types.IgniteDeploymentUpdateOpts{Name: "new name"},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		name,
		[]any{"test test", types.IgniteDeploymentUpdateOpts{Name: "new name"}, WithProjectID("test123")},
		&types.Deployment{ID: "hello"})
}

func TestClient_Ignite_Deployments_Patch(t *testing.T) {
	testIgniteUpdate(t, "Patch")
}

func TestClient_Ignite_Deployments_Update(t *testing.T) {
	testIgniteUpdate(t, "Update")
}

func TestClient_Ignite_Deployments_GetContainers(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/ignite/deployments/test%20test/containers",
		wantResultKey:  "containers",
		wantIgnore404:  false,
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetContainers",
		[]any{"test test", WithProjectID("test123")},
		[]*types.Container{{ID: "hello"}})
}

func TestClient_Ignite_Deployments_GetAll(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/ignite/deployments",
		wantResultKey:  "deployments",
		wantIgnore404:  false,
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetAll",
		[]any{WithProjectID("test123")},
		[]*types.Deployment{{ID: "hello"}})
}

func TestClient_Ignite_Deployments_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "DELETE",
		wantPath:       "/ignite/deployments/test%20test",
		wantIgnore404:  false,
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"Delete",
		[]any{"test test", WithProjectID("test123")},
		nil)
}

func TestClient_Ignite_Deployments_GetAllGateways(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "GET",
		wantPath:       "/ignite/deployments/test%20test/gateways",
		wantResultKey:  "gateways",
		wantIgnore404:  false,
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetAllGateways",
		[]any{"test test", WithProjectID("test123")},
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

func TestClient_Ignite_Deployments_Scale(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "PATCH",
		wantPath:       "/ignite/deployments/test%20test/scale",
		wantResultKey:  "containers",
		wantIgnore404:  false,
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		wantBody:       map[string]uint{"scale": 2},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"Scale",
		[]any{"test test", uint(2), WithProjectID("test123")},
		[]*types.Container{{ID: "hello"}})
}

func TestClient_Ignite_Containers_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "DELETE",
		wantPath:       "/ignite/containers/test%20test",
		wantIgnore404:  false,
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteContainers{c: c},
		"Delete",
		[]any{"test test", WithProjectID("test123")},
		nil)
}

func TestClient_Ignite_Containers_GetLogs(t *testing.T) {
	c := &mockClientDoer{}
	res := (&ClientCategoryIgniteContainers{c: c}).GetLogs(
		"test test", 10, true)
	assert.Equal(t, res, &Paginator[*types.ContainerLog]{
		c:           c,
		total:       -1,
		offsetStrat: true,
		limit:       10,
		path:        "/ignite/containers/test%20test/logs",
		resultKey:   "logs",
		sortBy:      "timestamp",
		orderBy:     "asc",
	})
}

func TestClient_Ignite_Containers_Stop(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "PUT",
		wantPath:       "/ignite/containers/test%20test/state",
		wantIgnore404:  false,
		wantBody:       map[string]types.ContainerState{"preferred_state": "stopped"},
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteContainers{c: c},
		"Stop",
		[]any{"test test", WithProjectID("test123")},
		nil)
}

func TestClient_Ignite_Containers_Start(t *testing.T) {
	c := &mockClientDoer{
		t:              t,
		wantMethod:     "PUT",
		wantPath:       "/ignite/containers/test%20test/state",
		wantIgnore404:  false,
		wantBody:       map[string]types.ContainerState{"preferred_state": "running"},
		wantClientOpts: []ClientOption{WithProjectID("test123")},
		tokenType:      "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteContainers{c: c},
		"Start",
		[]any{"test test", WithProjectID("test123")},
		nil)
}

func TestClient_Ignite_Deployments_NewHealthCheck(t *testing.T) {
	tests := []struct {
		name string

		input types.HealthCheckCreateOpts
		body  types.HealthCheckCreateOpts
	}{
		{
			name:  "defaults set",
			input: types.HealthCheckCreateOpts{DeploymentID: "test test"},
			body: types.HealthCheckCreateOpts{
				Protocol:     types.HealthCheckProtocolHTTP,
				Path:         "/",
				Port:         8080,
				InitialDelay: types.Seconds(time.Second * 5),
				Interval:     types.Seconds(time.Minute),
				Timeout:      types.Milliseconds(time.Millisecond * 50),
				MaxRetries:   3,
			},
		},
		{
			name: "all set by user",
			input: types.HealthCheckCreateOpts{
				DeploymentID: "test test",
				Protocol:     types.HealthCheckProtocolHTTP,
				Path:         "/hello",
				Port:         8000,
				InitialDelay: types.Seconds(time.Second * 10),
				Interval:     types.Seconds(time.Second * 5),
				Timeout:      types.Milliseconds(time.Millisecond * 150),
				MaxRetries:   10,
			},
			body: types.HealthCheckCreateOpts{
				Protocol:     types.HealthCheckProtocolHTTP,
				Path:         "/hello",
				Port:         8000,
				InitialDelay: types.Seconds(time.Second * 10),
				Interval:     types.Seconds(time.Second * 5),
				Timeout:      types.Milliseconds(time.Millisecond * 150),
				MaxRetries:   10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &mockClientDoer{
				t:             t,
				wantMethod:    "POST",
				wantPath:      "/ignite/deployments/test%20test/health-checks",
				wantResultKey: "health_check",
				wantIgnore404: false,
				tokenType:     "pat",
				wantBody:      tt.body,
			}
			testApiSingleton(c,
				&ClientCategoryIgniteDeployments{c: c},
				"NewHealthCheck",
				[]any{tt.input},
				&types.HealthCheck{ID: "hello"})
		})
	}
}

func TestClient_Ignite_Deployments_GetHealthChecks(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/deployments/test%20test/health-checks",
		wantResultKey: "health_checks",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetHealthChecks",
		[]any{"test test"},
		[]*types.HealthCheck{{ID: "hello"}})
}

func TestClient_Ignite_Deployments_DeleteHealthCheck(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "DELETE",
		wantPath:      "/ignite/deployments/test%20test/health-checks/testing%20123",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"DeleteHealthCheck",
		[]any{"test test", "testing 123"},
		nil)
}

func TestClient_Ignite_Deployments_UpdateHealthCheck(t *testing.T) {
	body := types.HealthCheckUpdateOpts{
		DeploymentID:  "test test",
		HealthCheckID: "hello world",
	}
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "PATCH",
		wantPath:      "/ignite/deployments/test%20test/health-checks/hello%20world",
		wantResultKey: "health_check",
		wantBody:      body,
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"UpdateHealthCheck",
		[]any{body},
		&types.HealthCheck{ID: "hello"})
}

func TestClient_Ignite_Deployments_HealthCheckStates(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/deployments/test%20test/health-check-state",
		wantResultKey: "health_check_states",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"HealthCheckStates",
		[]any{"test test"},
		[]*types.HealthCheckState{{HealthCheckID: "test"}})
}

func TestClient_Ignite_Deployments_GetStorageStats(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/ignite/deployments/test%20test/storage",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryIgniteDeployments{c: c},
		"GetStorageStats",
		[]any{"test test"},
		types.DeploymentStorageInfo{
			Volume: &types.DeploymentStorageSize{
				ProvisionedSize: 10,
				UsedSize:        20,
			},
		})
}
