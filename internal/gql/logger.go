package gql

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"github.com/graph-gophers/graphql-go/log"
	"go.uber.org/zap"
)

func Logger(log *zap.Logger) graphql.SchemaOpt {
	return graphql.Logger(&logger{Logger: log})
}

type logger struct {
	*zap.Logger
}

// LogPanic implements logger interface.
func (l *logger) LogPanic(ctx context.Context, value interface{}) {
	l.Error("gql panic", zap.Any("error", value))
}

var _ log.Logger = &logger{Logger: &zap.Logger{}}
