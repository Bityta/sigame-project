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
	
	defer func() {
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("REDIS_HOST")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Check defaults
	if cfg.Server.HTTPPort == "" {
		t.Error("HTTP port is empty")
	}

	if cfg.Server.WSPort == "" {
		t.Error("WS port is empty")
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
					HTTPPort: "8003",
					WSPort:   "8083",
					GRPCPort: "50053",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
			},
			wantErr: false,
		},
		{
			name: "missing HTTP port",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host: "localhost",
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

