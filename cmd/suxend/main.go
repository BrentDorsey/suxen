package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/travelaudience/suxen/internal/debug"
	"go.uber.org/zap"
)

var subsystem = "suxend"

func main() {
	var cfg configuration
	cfg.init()
	cfg.parse()

	log, err := NewLogger(subsystem, version, Config{
		Environment: cfg.logger.environment,
		Level:       cfg.logger.level,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := run(ctx, cfg, log); err != nil {
		log.Fatal("unexpected application termination", zap.Error(err))
	}
	log.Info("application termination")
}

func run(setupCtx context.Context, cfg configuration, log *zap.Logger) error {
	cnr := newContainer(log.Named("dependency-container"))
	err := cnr.init(setupCtx, cfg,
		cnr.initPrometheus,
		cnr.initCache,
		cnr.initHTTPClient,
		cnr.initNexusClient,
		cnr.initGraphQLQueryController,
		cnr.initHTTPHandlers,
	)
	if err != nil {
		return err
	}

	runCtx, cancel := context.WithCancel(context.Background()) // TODO: listen to signals
	defer cancel()

	go debug.Run(runCtx, log.Named("debug"), debug.Opts{
		Host:        cfg.host,
		Port:        cfg.port + 1,
		Environment: cfg.environment,
		Version:     version,
		Gatherer:    cnr.prometheus.gatherer,
	})

	log.Info("application started", zap.String("host", cfg.host), zap.Int("port", cfg.port), zap.String("version", version))

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.host, cfg.port), cnr.http.handler); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
