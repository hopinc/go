package hop

import (
	"context"
	"net/url"

	"go.hop.io/sdk/types"
)

// GetAll is used to get all images in a project.
func (c ClientCategoryRegistryImages) GetAll(ctx context.Context, opts ...ClientOption) ([]*types.Image, error) {
	var images []*types.Image
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/registry/images",
		ResultKey: "images",
		Result:    &images,
	}, opts)
	if err != nil {
		return nil, err
	}
	return images, nil
}

// GetManifest is used to get the manifest for an image.
func (c ClientCategoryRegistryImages) GetManifest(ctx context.Context, image string, opts ...ClientOption) ([]*types.ImageManifest, error) {
	var manifests []*types.ImageManifest
	err := c.c.do(ctx, ClientArgs{
		Method:    "GET",
		Path:      "/registry/images/" + url.PathEscape(image) + "/manifests",
		ResultKey: "manifest",
		Result:    &manifests,
	}, opts)
	if err != nil {
		return nil, err
	}
	return manifests, nil
}

// Delete is used to delete an image.
func (c ClientCategoryRegistryImages) Delete(ctx context.Context, image string, opts ...ClientOption) error {
	return c.c.do(ctx, ClientArgs{
		Method: "DELETE",
		Path:   "/registry/images/" + url.PathEscape(image),
	}, opts)
}
