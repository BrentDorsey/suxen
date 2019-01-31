package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	lru "github.com/hashicorp/golang-lru"
	"github.com/prometheus/client_golang/prometheus"
	promhttp "github.com/travelaudience/go-promhttp"
	"github.com/travelaudience/suxen/internal/gql"
	"github.com/travelaudience/suxen/internal/nexus"
	"go.uber.org/zap"
)

type initFunc func(context.Context, configuration) (io.Closer, error)

type container struct {
	log *zap.Logger

	closers    []io.Closer // TODO: use
	prometheus struct {
		registerer prometheus.Registerer
		gatherer   prometheus.Gatherer
	}
	http struct {
		client  *promhttp.Client
		handler http.Handler
	}
	nexus struct {
		client *nexus.Client
	}
	cache struct {
		manifest *lru.TwoQueueCache
	}
	graphql struct {
		controller struct {
			query *relay.Handler
		}
	}
}

func newContainer(log *zap.Logger) *container {
	return &container{
		log: log,
	}
}

func (c *container) init(ctx context.Context, cfg configuration, fns ...initFunc) error {
	for _, fn := range fns {
		closer, err := fn(ctx, cfg)
		if err != nil {
			return fmt.Errorf("dependency initialization failure: %s", err.Error())
		}

		c.closers = append(c.closers, closer) // TODO: use manager
	}

	return nil
}

func (c *container) initPrometheus(_ context.Context, _ configuration) (io.Closer, error) {
	c.prometheus.registerer = prometheus.DefaultRegisterer
	c.prometheus.gatherer = prometheus.DefaultGatherer

	return nil, nil
}

func (c *container) initHTTPClient(_ context.Context, cfg configuration) (io.Closer, error) {
	c.http.client = &promhttp.Client{
		Namespace: "suxen",
		Client: &http.Client{
			// CheckRedirect do not allow to follow redirects.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Registerer: c.prometheus.registerer,
	}
	return nil, nil
}

func (c *container) initCache(_ context.Context, cfg configuration) (io.Closer, error) {
	var err error
	c.cache.manifest, err = lru.New2Q(2000)
	if err != nil {
		return nil, fmt.Errorf("nexus client cache initialization failure: %s", err.Error())
	}
	return nil, nil
}

func (c *container) initNexusClient(_ context.Context, cfg configuration) (io.Closer, error) {
	client, err := c.http.client.ForRecipient(cfg.nexus.svc.address)
	if err != nil {
		return nil, fmt.Errorf("http client (for %s) cannot be initialized: %s", cfg.nexus.svc.address, err.Error())
	}

	c.nexus.client = &nexus.Client{
		Client:     client,
		Cache:      c.cache.manifest,
		Log:        c.log.Named("nexus-client"),
		Repository: cfg.nexus.repository,
		Address:    cfg.nexus.svc.address,
		AuthToken:  cfg.nexus.svc.authToken,
		SearchPath: cfg.nexus.searchPath,
	}

	return nil, nil
}

func (c *container) initGraphQLQueryController(_ context.Context, cfg configuration) (io.Closer, error) {
	schema, err := graphql.ParseSchema(graphQLSchema, &resolver{
		log:    c.log.Named("gql"),
		config: cfg,
		client: c.nexus.client,
	}, gql.Logger(c.log))
	if err != nil {
		return nil, fmt.Errorf("gql schema parsing failure: %s", err.Error())
	}

	c.graphql.controller.query = &relay.Handler{Schema: schema}
	return nil, nil
}

func (c *container) initHTTPHandlers(_ context.Context, cfg configuration) (io.Closer, error) {
	mux := &promhttp.ServeMux{
		ServeMux:  http.NewServeMux(),
		Namespace: "suxen",
	}
	mux.Handle("/explore", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(explorePage)
	}))
	mux.Handle("/", http.FileServer(http.Dir(cfg.static)))
	mux.Handle("/query", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, X-Requested-With")
		rw.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")

		if r.Method == http.MethodOptions {
			rw.Write([]byte{})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		c.graphql.controller.query.ServeHTTP(rw, r.WithContext(ctx))
	}))

	c.http.handler = mux

	return nil, c.prometheus.registerer.Register(mux)
}
