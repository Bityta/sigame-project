package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpass")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("PACK_SERVICE_HOST", "localhost")
	os.Setenv("PACK_SERVICE_PORT", "8084")

	defer func() {
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("PACK_SERVICE_HOST")
		os.Unsetenv("PACK_SERVICE_PORT")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.Server.HTTPPort == "" {
		t.Error("HTTP port is empty")
	}

	if cfg.Server.WSPort == "" {
		t.Error("WS port is empty")
	}

	if cfg.Database.Host == "" {
		t.Error("Database host is empty")
	}

	if cfg.Redis.Host == "" {
		t.Error("Redis host is empty")
	}

	if cfg.PackService.Host == "" {
		t.Error("Pack service host is empty")
	}
}

