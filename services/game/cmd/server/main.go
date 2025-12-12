package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	appGame "sigame/game/internal/application/game"
	grpcClient "sigame/game/internal/adapter/grpc/pack"
	"sigame/game/internal/infrastructure/config"
	"sigame/game/internal/infrastructure/logger"
	"sigame/game/internal/port"
	"sigame/game/internal/transport/ws"
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

	authClient, err := initAuthClient(cfg)
	if err != nil {
		logger.Warnf(nil, "Failed to connect to Auth Service: %v (continuing without token validation)", err)
	} else {
		defer authClient.Close()
		logger.Infof(nil, "Connected to Auth Service at %s", cfg.GetAuthServiceAddress())
	}

	hub := initWebSocketHub()
	go hub.Run()
	logger.Infof(nil, "WebSocket hub started")

	if err := restoreActiveGames(hub, packClient, repos); err != nil {
		logger.Warnf(nil, "Failed to restore active games: %v", err)
	}

	handlers := initHandlers(hub, packClient, repos, pgClient, redisClient)
	wsHandler := initWebSocketHandler(hub, authClient)
	router := initRouter(handlers, wsHandler)

	httpServer := createHTTPServer(cfg, router)

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infof(nil, "Shutting down Game Service...")

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	logger.Infof(nil, "Stopping WebSocket hub...")
	hub.Stop()
	logger.Infof(nil, "All game managers stopped")

	logger.Infof(nil, "Shutting down HTTP server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf(nil, "HTTP server shutdown error: %v", err)
	} else {
		logger.Infof(nil, "HTTP server stopped")
	}

	logger.Infof(nil, "Game Service stopped gracefully")
}

func restoreActiveGames(hub *ws.Hub, packClient *grpcClient.PackClient, repos *Repositories) error {
	ctx := context.Background()
	const maxActiveGames = 1000

	gameIDs, err := repos.RedisGameRepo.GetActiveGames(ctx, maxActiveGames)
	if err != nil {
		return fmt.Errorf("failed to get active games: %w", err)
	}

	if len(gameIDs) == 0 {
		logger.Infof(ctx, "No active games to restore")
		return nil
	}

	logger.Infof(ctx, "Restoring %d active games", len(gameIDs))

	restored := 0
	for _, gameID := range gameIDs {
		if err := restoreGame(ctx, gameID, hub, packClient, repos, repos.EventRepo); err != nil {
			logger.Errorf(ctx, "Failed to restore game %s: %v", gameID, err)
			continue
		}
		restored++
	}

	logger.Infof(ctx, "Restored %d/%d active games", restored, len(gameIDs))
	return nil
}

func restoreGame(ctx context.Context, gameID uuid.UUID, hub *ws.Hub, packClient *grpcClient.PackClient, repos *Repositories, eventLogger port.EventLogger) error {
	game, err := repos.RedisGameRepo.LoadGameState(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to load game state: %w", err)
	}

	if !game.Status.IsActive() {
		return nil
	}

	pack, err := packClient.GetPackContent(ctx, game.PackID)
	if err != nil {
		return fmt.Errorf("failed to load pack: %w", err)
	}

	manager := appGame.New(game, pack, hub, eventLogger, repos.GameRepo, repos.RedisGameRepo)
	hub.RegisterGameManager(gameID, manager)
	manager.Start()

	logger.Infof(ctx, "Restored game %s", gameID)
	return nil
}

