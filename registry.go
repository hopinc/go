package hopgo

import (
	"context"
	"net/url"

	"github.com/hopinc/hop-go/types"
)

// GetAll is used to get all images in a project.
func (c ClientCategoryRegistryImages) GetAll(ctx context.Context, projectId string) ([]*types.Image, error) {
	var images []*types.Image
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/registry/images",
		resultKey: "images",
		result:    &images,
		query:     getProjectIdParam(projectId),
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

// GetManifest is used to get the manifest for an image.
func (c ClientCategoryRegistryImages) GetManifest(ctx context.Context, projectId, image string) ([]*types.ImageManifest, error) {
	var manifests []*types.ImageManifest
	err := c.c.do(ctx, clientArgs{
		method:    "GET",
		path:      "/registry/images/" + url.PathEscape(image) + "/manifests",
		resultKey: "manifest",
		result:    &manifests,
		query:     getProjectIdParam(projectId),
	})
	if err != nil {
		return nil, err
	}
	return manifests, nil
}

// Delete is used to delete an image.
func (c ClientCategoryRegistryImages) Delete(ctx context.Context, projectId, image string) error {
	return c.c.do(ctx, clientArgs{
		method: "DELETE",
		path:   "/registry/images/" + url.PathEscape(image),
		query:  getProjectIdParam(projectId),
	})
}
