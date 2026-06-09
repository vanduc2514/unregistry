package containerd

import (
	"context"
	"fmt"

	"github.com/containerd/containerd/v2/client"
	"github.com/distribution/distribution/v3"
	"github.com/distribution/reference"
)

// repository implements distribution.Repository backed by the containerd content and image stores.
type repository struct {
	client         *client.Client
	name           reference.Named
	containerdName reference.Named
	blobStore      *blobStore
}

var _ distribution.Repository = &repository{}

func newRepository(ctx context.Context, client *client.Client, name reference.Named) (*repository, error) {
	containerdName, err := containerdImageName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &repository{
		client:         client,
		name:           name,
		containerdName: containerdName,
		blobStore: &blobStore{
			client: client,
			repo:   containerdName,
		},
	}, nil
}

// Named returns the name of the repository.
func (r *repository) Named() reference.Named {
	return r.name
}

// Manifests returns the manifest service for the repository backed by the containerd content store.
func (r *repository) Manifests(
	_ context.Context, _ ...distribution.ManifestServiceOption,
) (distribution.ManifestService, error) {
	return &manifestService{
		repo:      r.containerdName,
		blobStore: r.blobStore,
	}, nil
}

// Blobs returns the blob store for the repository backed by the containerd content store.
func (r *repository) Blobs(_ context.Context) distribution.BlobStore {
	return r.blobStore
}

// Tags returns the tag service for the repository backed by the containerd image store.
func (r *repository) Tags(_ context.Context) distribution.TagService {
	return &tagService{
		client: r.client,
		repo:   r.containerdName,
	}
}

func containerdImageName(ctx context.Context, name reference.Named) (reference.Named, error) {
	host, _ := ctx.Value("http.request.host").(string)
	if host == "" {
		return name, nil
	}

	return reference.WithName(fmt.Sprintf("%s/%s", host, name.Name()))
}
