package nexus

import (
	"time"

	"github.com/Masterminds/semver"
	"github.com/graph-gophers/graphql-go"
)

type SearchResponse struct {
	Items             []*Item `json:"items"`
	ContinuationToken string
}

type Item struct {
	ID         string
	Repository string
	Format     string
	Group      string
	Name       string
	Version    string
	Assets     []Asset
}

type Asset struct {
	ID          string   `json:"id"`
	DownloadURL string   `json:"downloadUrl"`
	Path        string   `json:"path"`
	Checksum    Checksum `json:"checksum"`
}

type Checksum struct {
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
}

type Image struct {
	Name        string
	Tag         string
	ID          graphql.ID
	ManifestURL string
	Path        string
	Version     string
	Author      string
	CreatedAt   time.Time
}

type Manifest struct {
	Name            string
	Tag             string
	CreatedAt       time.Time `json:"created"`
	DockerVersion   string    `json:"docker_version"`
	OperatingSystem string    `json:"os"`
	Author          string    `json:"author"`
}

type ImageGroup struct {
	ID          graphql.ID
	CreatedAt   time.Time
	Name        string
	Version     *semver.Version
	Images      []*Image
	Checksum    Checksum
	DownloadURL string
	Manifest    *Manifest
}

type ByCreatedAt []ImageGroup

func (a ByCreatedAt) Len() int { return len(a) }
func (a ByCreatedAt) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByCreatedAt) Less(i, j int) bool {
	return a[i].CreatedAt.After(a[j].CreatedAt)
}
