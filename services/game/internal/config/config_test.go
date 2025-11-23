package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set required environment variables
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("REDIS_URL")
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig() returned nil config")
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
					URL: "postgres://test:test@localhost:5432/test",
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

