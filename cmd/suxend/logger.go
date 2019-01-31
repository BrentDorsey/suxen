package main

import (
	"errors"

	"github.com/piotrkowalczuk/zapstackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	production  = "production"
	development = "development"
	stackdriver = "stackdriver"
)

// NewLogger  allocates new logger based on given options.
func NewLogger(service, version string, pkgCfg Config) (logger *zap.Logger, err error) {
	if service == "" {
		return nil, errors.New("nexusuid: service name is missing, logger cannot be initialized")
	}
	if version == "" {
		return nil, errors.New("nexusuid: service version is missing, logger cannot be initialized")
	}

	var (
		zapCfg  zap.Config
		zapOpts []zap.Option
		lvl     zapcore.Level
	)
	switch pkgCfg.Environment {
	case production:
		zapCfg = zap.NewProductionConfig()
	case stackdriver:
		zapCfg = zapstackdriver.NewStackdriverConfig()
	case development:
		zapCfg = zap.NewDevelopmentConfig()
	default:
		zapCfg = zap.NewProductionConfig()
	}

	if err = lvl.Set(pkgCfg.Level); err != nil {
		return nil, err
	}
	zapCfg.Level.SetLevel(lvl)

	logger, err = zapCfg.Build(zapOpts...)
	if err != nil {
		return nil, err
	}
	logger = logger.With(zap.Object("serviceContext", &zapstackdriver.ServiceContext{
		Service: service,
		Version: version,
	}))
	logger.Info("logger has been initialized", zap.String("environment", pkgCfg.Environment))

	return logger, nil
}

type Config struct {
	Environment string
	Level       string
}
