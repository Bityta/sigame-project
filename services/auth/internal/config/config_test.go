package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set required environment variables
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpass")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("JWT_SECRET", "test-secret-key-min-32-chars-long")
	
	defer func() {
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.JWT.Secret == "" {
		t.Error("JWT secret is empty")
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected database host localhost, got %s", cfg.Database.Host)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8080",
					GRPCPort: "50051",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				JWT: JWTConfig{
					Secret: "test-secret-key-min-32-chars-long",
				},
			},
			wantErr: false,
		},
		{
			name: "missing JWT secret",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8080",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				JWT: JWTConfig{
					Secret: "",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database host",
			config: Config{
				JWT: JWTConfig{
					Secret: "test-secret",
				},
				Database: DatabaseConfig{
					Host: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_GetPostgresConnectionString(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	connStr := cfg.GetPostgresConnectionString()
	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"

	if connStr != expected {
		t.Errorf("Expected %s, got %s", expected, connStr)
	}
}

