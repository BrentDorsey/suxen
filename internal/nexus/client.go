package nexus

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver"
	"github.com/graph-gophers/graphql-go"
	"github.com/hashicorp/golang-lru"
	"go.uber.org/zap"
)

const headerContentType = "Content-Type"

type Client struct {
	*http.Client
	Log                                        *zap.Logger
	Cache                                      *lru.TwoQueueCache
	Address, SearchPath, Repository, AuthToken string
}

func (c *Client) get(ctx context.Context, uri string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	if c.AuthToken != "" {
		req.Header.Add("Authorization", "Basic "+c.AuthToken)
	}

	c.Log.Debug("outgoing nexus get request started", zap.String("uri", uri))

	start := time.Now()
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		c.Log.Debug("outgoing nexus get request finished (failure)",
			zap.String("uri", uri),
			zap.Duration("elapsed", time.Since(start)),
			zap.Error(err),
		)
		return nil, err
	}

	c.Log.Debug("outgoing nexus get request finished (completed)", zap.String("uri", uri),
		zap.Int("status_code", res.StatusCode),
		zap.String("content_type", res.Header.Get(headerContentType)),
		zap.Duration("elapsed", time.Since(start)),
	)
	return res, nil
}

func (c *Client) isResponseValid(res *http.Response, contentType string) error {
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status code: %s", http.StatusText(res.StatusCode))
	}

	if res.Header.Get(headerContentType) != contentType {
		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("wrong content type, got %s expected %s:\n%s", res.Header.Get(headerContentType), contentType, string(buf))
	}
	return nil
}

func (c *Client) Search(ctx context.Context, query string) ([]*Item, error) {
	var (
		continuationToken string
		assets            []*Item
	)

	for {
		res, err := c.search(ctx, continuationToken)
		if err != nil {
			return nil, err
		}

	SearchLoop:
		for _, item := range res.Items {
			if strings.Contains(item.Name, query) {
				assets = append(assets, item)
				continue
			}
			if strings.Contains(item.Version, query) {
				assets = append(assets, item)
				continue
			}
			if strings.Contains(item.ID, query) {
				assets = append(assets, item)
				continue
			}
			for _, asset := range item.Assets {
				if strings.Contains(asset.Path, query) {
					assets = append(assets, item)
					continue SearchLoop
				}
				if strings.Contains(asset.DownloadURL, query) {
					assets = append(assets, item)
					continue SearchLoop
				}
			}
		}

		if res.ContinuationToken == "" {
			break
		}

		continuationToken = res.ContinuationToken
		// TODO: implement infinite scroll
		break
	}

	return assets, nil
}

func (c *Client) search(ctx context.Context, continuationToken string) (*SearchResponse, error) {
	const contentType = "application/json"

	uri := fmt.Sprintf("%s/%s?repository=%s", c.Address, c.SearchPath, c.Repository)
	if continuationToken != "" {
		uri = fmt.Sprintf("%s&continuationToken=%s", uri, continuationToken)
	}

	c.Log.Debug("call to nexus", zap.String("uri", uri))

	res, err := c.get(ctx, uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := c.isResponseValid(res, contentType); err != nil {
		return nil, err
	}

	var data SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("cannot unarshal search response due to error: %s", err.Error())
	}

	c.Log.Debug("response from nexus", zap.Int("count", len(data.Items)))

	return &data, nil
}

func (c *Client) IntoImageGroups(ctx context.Context, items []*Item) []ImageGroup {
	grouping := make(map[Checksum][]*Item)

	for _, item := range items {
		for _, asset := range item.Assets {
			if _, ok := grouping[asset.Checksum]; !ok {
				grouping[asset.Checksum] = []*Item{}
			}
			grouping[asset.Checksum] = append(grouping[asset.Checksum], item)
		}
	}

	pipe := make(chan ImageGroup, len(grouping))
	var wg sync.WaitGroup
	go func() {
		// sem keeps number of outgoing request fixed for a single incoming request.
		sem := make(chan struct{}, 5)
		for _, group := range grouping {
			wg.Add(1)

			go func(group []*Item) {
				defer wg.Done()

				var (
					ig ImageGroup
				)
				for _, item := range group {
					var majorImage bool
					if ig.Name == "" {
						ig.Name = item.Name
					}

					if ver, err := semver.NewVersion(item.Version); err == nil {
						if ig.Version == nil {
							ig.Version = ver
							majorImage = true
						} else {
							if ig.Version.LessThan(ver) {
								ig.Version = ver
								majorImage = true
							}
						}
					}

					for i, asset := range item.Assets {
						img := &Image{
							ID:          graphql.ID(asset.ID),
							ManifestURL: asset.DownloadURL,
							Path:        asset.Path,
							Version:     item.Version,
						}
						if asset.DownloadURL != "" {
							var (
								manifest *Manifest
								err      error
							)
							if val, ok := c.Cache.Get(asset.DownloadURL); ok {
								manifest = val.(*Manifest)
							} else {
								sem <- struct{}{}
								manifest, err = c.fetchManifest(ctx, asset.DownloadURL)
								if err != nil {
									c.Log.Debug("manifest fetch failure", zap.Error(err), zap.String("download_url", asset.DownloadURL))

									switch err {
									case context.DeadlineExceeded, context.Canceled:
										<-sem
										return
									}
								} else {
									c.Cache.Add(asset.DownloadURL, manifest)
								}
								<-sem
							}
							if err == nil {
								ig.CreatedAt = manifest.CreatedAt
								ig.Manifest = manifest
								img.Name = manifest.Name
								img.Tag = manifest.Tag
								img.Author = manifest.Author
								img.CreatedAt = manifest.CreatedAt
							}
						}

						ig.Checksum = asset.Checksum

						if ig.ID == "" {
							ig.ID = graphql.ID(asset.Checksum.Sha256)
						}
						if i == 0 && majorImage {
							ig.ID = graphql.ID(asset.Checksum.Sha256)
							ig.DownloadURL = asset.DownloadURL
						}
						ig.Images = append(ig.Images, img)
					}
				}
				pipe <- ig
			}(group)
		}
		wg.Wait()
		close(pipe)
	}()

	var res []ImageGroup
	for ig := range pipe {
		res = append(res, ig)
	}

	sort.Sort(ByCreatedAt(res))

	return res
}

func (c *Client) fetchManifest(ctx context.Context, downloadURL string) (*Manifest, error) {
	const contentType = "application/vnd.docker.distribution.manifest.v1+json"

	res, err := c.get(ctx, downloadURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := c.isResponseValid(res, contentType); err != nil {
		return nil, err
	}

	var data struct {
		Name    string
		Tag     string
		History []struct {
			V1Compatibility string `json:"v1Compatibility"`
		} `json:"history"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("cannot unarshal manifest (%s) due to error: %s", data.Name, err.Error())
	}
	if len(data.History) < 1 {
		return nil, nil
	}

	var manifest Manifest
	if err := json.Unmarshal([]byte(data.History[0].V1Compatibility), &manifest); err != nil {
		return nil, fmt.Errorf("cannot unarshal field history of manifest (%s) due to error: %s", data.Name, err.Error())
	}

	manifest.Name = data.Name
	manifest.Tag = data.Tag

	return &manifest, nil
}
