package hop

import (
	"testing"

	"github.com/hopinc/hop-go/types"
)

func TestClient_Registry_GetAll(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/registry/images",
		wantQuery:     map[string]string{"project": "test123"},
		wantResultKey: "images",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryRegistryImages{c: c},
		"GetAll",
		[]any{"test123"},
		[]*types.Image{{Name: "hello"}})
}

func TestClient_Registry_GetManifest(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "GET",
		wantPath:      "/registry/images/test%20test/manifests",
		wantQuery:     map[string]string{"project": "test123"},
		wantResultKey: "manifest",
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryRegistryImages{c: c},
		"GetManifest",
		[]any{"test123", "test test"},
		[]*types.ImageManifest{{Tag: types.StringPointerify("hello")}})
}

func TestClient_Registry_Delete(t *testing.T) {
	c := &mockClientDoer{
		t:             t,
		wantMethod:    "DELETE",
		wantPath:      "/registry/images/test%20test",
		wantQuery:     map[string]string{"project": "test123"},
		wantIgnore404: false,
		tokenType:     "pat",
	}
	testApiSingleton(c,
		&ClientCategoryRegistryImages{c: c},
		"Delete",
		[]any{"test123", "test test"},
		nil)
}
