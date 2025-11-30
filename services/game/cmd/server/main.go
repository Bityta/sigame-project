package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	
	"github.com/sigame/game/internal/config"
	grpcClient "github.com/sigame/game/internal/grpc"
	"github.com/sigame/game/internal/metrics"
	"github.com/sigame/game/internal/repository/postgres"
	"github.com/sigame/game/internal/repository/redis"
	"github.com/sigame/game/internal/tracing"
	"github.com/sigame/game/internal/transport/rest"
	"github.com/sigame/game/internal/transport/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting Game Service...")
	log.Printf("HTTP Port: %s", cfg.Server.HTTPPort)
	log.Printf("WS Port: %s", cfg.Server.WSPort)
	log.Printf("Mode: %s", cfg.Server.Mode)

	// Initialize OpenTelemetry tracer
	tp, err := tracing.InitTracer("game-service")
	if err != nil {
		log.Printf("Warning: Failed to initialize tracer: %v", err)
	} else {
		defer tracing.Shutdown(tp)
	}

	// Initialize metrics
	_ = metrics.NewMetrics()
	log.Printf("✓ Metrics initialized")

	// Connect to PostgreSQL
	pgClient, err := postgres.NewClient(
		cfg.GetPostgresConnectionString(),
		cfg.Database.MaxConns,
		cfg.Database.MaxIdle,
	)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgClient.Close()
	log.Printf("✓ Connected to PostgreSQL")

	// Create repositories
	gameRepo := postgres.NewGameRepository(pgClient.GetDB())
	eventRepo := postgres.NewEventRepository(pgClient.GetDB())

	// Connect to Redis
	redisClient, err := redis.NewClient(
		cfg.GetRedisAddress(),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Printf("✓ Connected to Redis")

	// Create Redis repositories
	redisGameRepo := redis.NewGameRepository(redisClient.GetClient())
	redisCacheRepo := redis.NewCacheRepository(redisClient.GetClient())

	// Connect to Pack Service (gRPC)
	packClient, err := grpcClient.NewPackClient(cfg.GetPackServiceAddress())
	if err != nil {
		log.Fatalf("Failed to connect to Pack Service: %v", err)
	}
	defer packClient.Close()
	log.Printf("✓ Connected to Pack Service at %s", cfg.GetPackServiceAddress())

	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()
	log.Printf("✓ WebSocket hub started")

	// Create handlers
	wsHandler := websocket.NewHandler(hub)
	restHandler := rest.NewHandler(hub, packClient, gameRepo, eventRepo, redisGameRepo, redisCacheRepo)

	// Setup router
	router := rest.SetupRouter(restHandler, wsHandler)
	
	// Add OpenTelemetry middleware
	router.Use(otelgin.Middleware("game-service"))

	// Create HTTP server
	httpAddr := fmt.Sprintf(":%s", cfg.Server.HTTPPort)
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	// Start metrics server
	metricsAddr := ":9090"
	metricsServer := &http.Server{
		Addr:    metricsAddr,
		Handler: promhttp.Handler(),
	}

	go func() {
		log.Printf("✓ Metrics server listening on %s", metricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	// Start HTTP server
	go func() {
		log.Printf("✓ HTTP server listening on %s", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	log.Printf("==============================================")
	log.Printf("Game Service is ready!")
	log.Printf("HTTP API: http://localhost%s", httpAddr)
	log.Printf("Metrics: http://localhost%s/metrics", metricsAddr)
	log.Printf("==============================================")

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Game Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Shutdown metrics server
	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Printf("Metrics server shutdown error: %v", err)
	}

	log.Println("Game Service stopped gracefully")
}

