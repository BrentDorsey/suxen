package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Opts struct {
	Host        string
	Port        int
	Environment string
	Version     string
	Gatherer    prometheus.Gatherer
}

func Run(ctx context.Context, log *zap.Logger, opts Opts) {
	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)

	log.Info("debug server is running", zap.String("address", addr))

	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	mux.Handle("/healthz", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(rw).Encode(struct {
			Version     string `json:"version"`
			Environment string `json:"environment"`
			DebugAddr   string `json:"debugAddress"`
		}{
			Version:     opts.Version,
			Environment: opts.Environment,
			DebugAddr:   addr,
		}); err != nil {
			log.Error("healthz info cannot be encoded", zap.Error(err))
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	}))
	mux.Handle("/healthr", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		select {
		case <-ctx.Done():
			http.Error(rw, ctx.Err().Error(), http.StatusServiceUnavailable)
		default:
			io.WriteString(rw, "1")
			//rw.WriteHeader(http.StatusOK)
		}
	}))
	mux.Handle("/metrics", promhttp.HandlerFor(opts.Gatherer, promhttp.HandlerOpts{
		ErrorLog: zap.NewStdLog(log),
	}))

	err := http.ListenAndServe(addr, mux)
	log.Error("debug server failure", zap.Error(err))
}
