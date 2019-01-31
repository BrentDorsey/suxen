package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/travelaudience/suxen/internal/nexus"
	"go.uber.org/zap"
)

type resolver struct {
	log    *zap.Logger
	client *nexus.Client
	config configuration
}

func (r *resolver) Search(ctx context.Context, args struct{ Query *string }) ([]*imageGroupResolver, error) {
	var query string
	if args.Query != nil {
		query = *args.Query
	}

	r.log.Debug("incoming search request", zap.String("query", query))

	items, err := r.client.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	var res []*imageGroupResolver
	for _, group := range r.client.IntoImageGroups(ctx, items) {
		res = append(res, &imageGroupResolver{val: group, config: r.config})
	}
	return res, nil
}

type imageGroupResolver struct {
	config configuration
	val    nexus.ImageGroup
}

func (r *imageGroupResolver) PreRelease() (*string, error) {
	if r.val.Version == nil {
		return nil, nil
	}
	pr := r.val.Version.Prerelease()
	return &pr, nil
}

func (r *imageGroupResolver) ID() *graphql.ID {
	return &r.val.ID
}

func (r *imageGroupResolver) Name() string {
	return r.val.Name
}

func (r *imageGroupResolver) Version() string {
	if r.val.Version == nil {
		return ""
	}
	return r.val.Version.Original()
}

func (r *imageGroupResolver) Images() []*imageResolver {
	var res []*imageResolver
	for _, img := range r.val.Images {
		res = append(res, &imageResolver{val: img, config: r.config})
	}
	return res
}

func (r *imageGroupResolver) Checksum() *checksumResolver {
	return &checksumResolver{val: r.val.Checksum}
}

func (r *imageGroupResolver) CreatedAt() *string {
	res := r.val.CreatedAt.Format(time.RFC3339)
	return &res
}

func (r *imageGroupResolver) DockerVersion() *string {
	return &r.val.Manifest.DockerVersion
}

func (r *imageGroupResolver) OperatingSystem() *string {
	return &r.val.Manifest.OperatingSystem
}

func (r *imageGroupResolver) Author() *string {
	return &r.val.Manifest.Author
}

type imageResolver struct {
	config configuration
	val    *nexus.Image
}

func (r *imageResolver) ID() *graphql.ID {
	return &r.val.ID
}

func (r *imageResolver) Name() string {
	return r.val.Name
}

func (r *imageResolver) Version() string {
	return r.val.Tag
}

func (r *imageResolver) ManifestUrl() string {
	return strings.Replace(r.val.ManifestURL, r.config.nexus.svc.address, r.config.nexus.address, 1)
}

func (r *imageResolver) PullUrl() string {
	return fmt.Sprintf("%s/%s:%s", r.config.nexus.registry.address, r.val.Name, r.val.Tag)
}

func (r *imageResolver) Path() string {
	return r.val.Path
}

func (r *imageResolver) Author() *string {
	return &r.val.Author
}

func (r *imageResolver) CreatedAt() *string {
	res := r.val.CreatedAt.Format(time.RFC3339)
	return &res
}

type checksumResolver struct {
	val nexus.Checksum
}

func (r *checksumResolver) Sha1() string {
	return r.val.Sha1
}

func (r *checksumResolver) Sha256() string {
	return r.val.Sha256
}
