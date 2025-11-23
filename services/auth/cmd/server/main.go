package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	
	"github.com/sigame/auth/internal/config"
	"github.com/sigame/auth/internal/logger"
	"github.com/sigame/auth/internal/metrics"
	"github.com/sigame/auth/internal/repository/postgres"
	redisrepo "github.com/sigame/auth/internal/repository/redis"
	"github.com/sigame/auth/internal/service"
	"github.com/sigame/auth/internal/tracing"
	grpcTransport "github.com/sigame/auth/internal/transport/grpc"
	"github.com/sigame/auth/internal/transport/rest"
	pb "github.com/sigame/auth/proto"
)

func main() {
	ctx := context.Background()
	
	// Initialize logger
	logger.Init("auth-service")
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Errorf(ctx, "Failed to load configuration: %v", err)
		os.Exit(1)
	}

	logger.Info(ctx, "Starting Auth Service...")
	logger.Infof(ctx, "HTTP Server: :%s", cfg.Server.HTTPPort)
	logger.Infof(ctx, "gRPC Server: :%s", cfg.Server.GRPCPort)

	// Initialize OpenTelemetry tracer
	tp, err := initTracer(ctx)
	if err != nil {
		logger.Warnf(ctx, "Failed to initialize tracer: %v", err)
	} else {
		defer tracing.Shutdown(tp)
	}

	// Initialize database connections
	db, redisClient := initDatabases(ctx, cfg)
	defer db.Close()
	defer redisClient.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	cacheRepo := redisrepo.NewCacheRepository(redisClient)

	// Initialize services
	authService := initServices(userRepo, cacheRepo, cfg)

	// Initialize metrics
	metricsRegistry := metrics.New()

	// Setup HTTP server
	httpServer := setupHTTPServer(ctx, cfg, authService, metricsRegistry)

	// Setup gRPC server
	grpcServer := setupGRPCServer(authService, metricsRegistry)

	// Start servers
	startServers(ctx, cfg, httpServer, grpcServer)

	// Start metrics updater
	startMetricsUpdater(ctx, userRepo, cacheRepo, metricsRegistry)

	// Start token cleanup job
	startTokenCleanupJob(ctx, userRepo)

	logger.Info(ctx, "✓ Auth Service started successfully")
	logger.Info(ctx, "✓ Ready to accept requests")

	// Wait for shutdown signal and handle graceful shutdown
	handleShutdown(ctx, httpServer, grpcServer)
}

// initTracer initializes the OpenTelemetry tracer
func initTracer(ctx context.Context) (*tracing.TracerProvider, error) {
	return tracing.InitTracer("auth-service")
}

// initDatabases initializes PostgreSQL and Redis connections
func initDatabases(ctx context.Context, cfg *config.Config) (*sql.DB, *redis.Client) {
	// Connect to PostgreSQL
	logger.Info(ctx, "Connecting to PostgreSQL...")
	db, err := postgres.Connect(
		cfg.GetPostgresConnectionString(),
		cfg.Database.MaxConns,
		cfg.Database.MaxIdle,
	)
	if err != nil {
		logger.Errorf(ctx, "Failed to connect to PostgreSQL: %v", err)
		os.Exit(1)
	}
	logger.Info(ctx, "✓ Connected to PostgreSQL")

	// Connect to Redis
	logger.Info(ctx, "Connecting to Redis...")
	redisClient, err := redisrepo.Connect(
		cfg.GetRedisAddress(),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		logger.Errorf(ctx, "Failed to connect to Redis: %v", err)
		os.Exit(1)
	}
	logger.Info(ctx, "✓ Connected to Redis")

	return db, redisClient
}

// initServices initializes the auth and JWT services
func initServices(userRepo *postgres.UserRepository, cacheRepo *redisrepo.CacheRepository, cfg *config.Config) *service.AuthService {
	// Initialize JWT service
	jwtService := service.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// Initialize auth service
	return service.NewAuthService(
		userRepo,
		cacheRepo,
		jwtService,
		service.RateLimitConfig{
			Attempts: cfg.RateLimit.Attempts,
			Window:   cfg.RateLimit.Window,
		},
	)
}

// setupHTTPServer sets up the HTTP server with all middleware and routes
func setupHTTPServer(ctx context.Context, cfg *config.Config, authService *service.AuthService, metricsRegistry *metrics.Metrics) *http.Server {
	// Initialize HTTP handlers
	httpHandler := rest.NewHandler(authService)
	jwtMiddleware := rest.JWTAuthMiddleware(authService)

	// Setup HTTP router
	router := rest.SetupRouter(httpHandler, jwtMiddleware, metricsRegistry)
	
	// Add OpenTelemetry middleware
	router.Use(otelgin.Middleware("auth-service"))

	// Setup metrics endpoint
	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	return &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// setupGRPCServer sets up the gRPC server with interceptors
func setupGRPCServer(authService *service.AuthService, metricsRegistry *metrics.Metrics) *grpc.Server {
	// Initialize gRPC server with metrics and tracing
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcTransport.MetricsInterceptor(metricsRegistry)),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	authGRPCServer := grpcTransport.NewServer(authService)
	pb.RegisterAuthServiceServer(grpcServer, authGRPCServer)
	
	// Register reflection service for grpcurl support
	reflection.Register(grpcServer)

	return grpcServer
}

// startServers starts the HTTP and gRPC servers in goroutines
func startServers(ctx context.Context, cfg *config.Config, httpServer *http.Server, grpcServer *grpc.Server) {
	// Start HTTP server in goroutine
	go func() {
		logger.Infof(ctx, "HTTP server listening on :%s", cfg.Server.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf(ctx, "HTTP server failed: %v", err)
			os.Exit(1)
		}
	}()

	// Start gRPC server in goroutine
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
		if err != nil {
			logger.Errorf(ctx, "Failed to listen on gRPC port: %v", err)
			os.Exit(1)
		}
		logger.Infof(ctx, "gRPC server listening on :%s", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Errorf(ctx, "gRPC server failed: %v", err)
			os.Exit(1)
		}
	}()
}

// startMetricsUpdater starts a background goroutine to update metrics periodically
func startMetricsUpdater(ctx context.Context, userRepo *postgres.UserRepository, cacheRepo *redisrepo.CacheRepository, metricsRegistry *metrics.Metrics) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		// Update metrics immediately on startup
		updateMetrics(ctx, userRepo, cacheRepo, metricsRegistry)
		
		for range ticker.C {
			updateMetrics(ctx, userRepo, cacheRepo, metricsRegistry)
		}
	}()
}

// startTokenCleanupJob starts a background goroutine to clean up expired tokens periodically
func startTokenCleanupJob(ctx context.Context, userRepo *postgres.UserRepository) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		
		for range ticker.C {
			deleted, err := userRepo.DeleteExpiredTokens(ctx)
			if err != nil {
				logger.Errorf(ctx, "Failed to delete expired tokens: %v", err)
			} else if deleted > 0 {
				logger.Infof(ctx, "Deleted %d expired refresh tokens", deleted)
			}
		}
	}()
}

// handleShutdown waits for shutdown signal and performs graceful shutdown
func handleShutdown(ctx context.Context, httpServer *http.Server, grpcServer *grpc.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "Shutting down servers...")

	// Shutdown gRPC server gracefully
	logger.Info(ctx, "Stopping gRPC server...")
	grpcServer.GracefulStop()

	// Shutdown HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Errorf(ctx, "HTTP server forced to shutdown: %v", err)
	}

	logger.Info(ctx, "✓ Auth Service stopped")
}

// updateMetrics updates user and session metrics from database
func updateMetrics(ctx context.Context, userRepo *postgres.UserRepository, cacheRepo *redisrepo.CacheRepository, m *metrics.Metrics) {
	// Count total users
	totalUsers, err := userRepo.CountUsers(ctx)
	if err != nil {
		logger.Errorf(ctx, "Failed to count users: %v", err)
	} else {
		m.SetTotalUsers(totalUsers)
	}
	
	// Count active sessions from Redis
	activeSessions, err := cacheRepo.CountActiveSessions(ctx)
	if err != nil {
		logger.Errorf(ctx, "Failed to count active sessions: %v", err)
	} else {
		m.SetActiveSessions(activeSessions)
	}
}

