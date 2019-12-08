package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type MetricsServer struct {
	httpServer *http.Server
}

func (s *MetricsServer) Serve() {
	handler := http.NewServeMux()
	s.httpServer = &http.Server{
		Addr:              ":2112",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	handler.Handle("/metrics", promhttp.Handler())
	err := s.httpServer.ListenAndServe()
	log.Warningf("http metrics server exited with %v", err)
}

func (s *MetricsServer) Close() {
	log.Debug("closing metrics http server")
	_ = s.httpServer.Close()
}
