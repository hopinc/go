package types

// ImageDigest is used to define the digest of an image.
type ImageDigest struct {
	Digest   string `json:"digest"`
	Size     int    `json:"size"`
	Uploaded string `json:"uploaded"`
}

// ImageManifest is used to define the manifest for an image.
type ImageManifest struct {
	// Digest is the digest of the image.
	Digest ImageDigest `json:"digest"`

	// Tag is used to define an image tag.
	Tag *string `json:"tag"`
}
