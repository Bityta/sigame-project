package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"sigame/game/internal/infrastructure/config"
	"sigame/game/internal/infrastructure/logger"
	"sigame/game/internal/infrastructure/tracing"
)

func main() {
	logger.Init(ServiceName)

	cfg, err := config.Load()
	if err != nil {
		logger.Errorf(nil, "Failed to load config: %v", err)
		os.Exit(1)
	}

	logger.Infof(nil, "Starting Game Service...")
	logger.Infof(nil, "HTTP Port: %s", cfg.Server.HTTPPort)
	logger.Infof(nil, "WS Port: %s", cfg.Server.WSPort)

	tp, err := initTracer(ServiceName)
	if err != nil {
		logger.Warnf(nil, "Failed to initialize tracer: %v", err)
	} else if tp != nil {
		defer tracing.Shutdown(tp)
	}

	_ = initMetrics()
	logger.Infof(nil, "Metrics initialized")

	pgClient, err := initPostgreSQL(cfg)
	if err != nil {
		logger.Errorf(nil, "Failed to connect to PostgreSQL: %v", err)
		os.Exit(1)
	}
	defer pgClient.Close()
	logger.Infof(nil, "Connected to PostgreSQL")

	redisClient, err := initRedis(cfg)
	if err != nil {
		logger.Errorf(nil, "Failed to connect to Redis: %v", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	logger.Infof(nil, "Connected to Redis")

	repos := initRepositories(pgClient, redisClient)

	packClient, err := initPackClient(cfg)
	if err != nil {
		logger.Errorf(nil, "Failed to connect to Pack Service: %v", err)
		os.Exit(1)
	}
	defer packClient.Close()
	logger.Infof(nil, "Connected to Pack Service at %s", cfg.GetPackServiceAddress())

	hub := initWebSocketHub()
	go hub.Run()
	logger.Infof(nil, "WebSocket hub started")

	handlers := initHandlers(hub, packClient, repos)
	wsHandler := initWebSocketHandler(hub)
	router := initRouter(handlers, wsHandler)

	httpServer := createHTTPServer(cfg, router)
	metricsServer := createMetricsServer()

	go func() {
		logger.Infof(nil, "Metrics server listening on %s", MetricsPort)
		if err := startMetricsServer(metricsServer, MetricsPort); err != nil && err != http.ErrServerClosed {
			logger.Errorf(nil, "Metrics server error: %v", err)
		}
	}()

	go func() {
		httpAddr := fmt.Sprintf(":%s", cfg.Server.HTTPPort)
		logger.Infof(nil, "HTTP server listening on %s", httpAddr)
		if err := startHTTPServer(httpServer, httpAddr); err != nil && err != http.ErrServerClosed {
			logger.Errorf(nil, "HTTP server error: %v", err)
			os.Exit(1)
		}
	}()

	httpAddr := fmt.Sprintf(":%s", cfg.Server.HTTPPort)
	logger.Infof(nil, "Game Service is ready!")
	logger.Infof(nil, "HTTP API: %s", fmt.Sprintf("http://localhost%s", httpAddr))
	logger.Infof(nil, "Metrics: %s", fmt.Sprintf("http://localhost%s/metrics", MetricsPort))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infof(nil, "Shutting down Game Service...")

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := shutdownServers(ctx, httpServer, metricsServer); err != nil {
		logger.Errorf(nil, "Shutdown error: %v", err)
	}

	logger.Infof(nil, "Game Service stopped gracefully")
}

