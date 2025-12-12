package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sigame/game/internal/infrastructure/config"
)

func createHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	httpAddr := fmt.Sprintf(":%s", cfg.Server.HTTPPort)
	return &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}
}

func createMetricsServer() *http.Server {
	return &http.Server{
		Addr:    MetricsPort,
		Handler: promhttp.Handler(),
	}
}

func startHTTPServer(server *http.Server, addr string) error {
	return server.ListenAndServe()
}

func startMetricsServer(server *http.Server, addr string) error {
	return server.ListenAndServe()
}

func shutdownServers(ctx context.Context, httpServer, metricsServer *http.Server) error {
	var errs []error

	if err := httpServer.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("HTTP server shutdown error: %w", err))
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("metrics server shutdown error: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

