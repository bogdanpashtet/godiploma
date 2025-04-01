package app

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/bogdanpashtet/godiploma/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// newMetricsServer returns a new HTTP server that serves metrics.
func newMetricsServer(cfg config.MetricsConfig) *http.Server {
	mux := http.NewServeMux()

	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = "/metrics"
	}
	mux.HandleFunc(endpoint, promhttp.Handler().ServeHTTP)

	attachPprof(mux)

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
}

func attachPprof(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
