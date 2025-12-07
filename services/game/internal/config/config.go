package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the game service
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	PackService PackServiceConfig
}

// ServerConfig holds HTTP, WebSocket and gRPC server configuration
type ServerConfig struct {
	HTTPPort string
	WSPort   string
	GRPCPort string
	Mode     string // "debug" or "release"
}

// DatabaseConfig holds PostgreSQL database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MaxIdle  int
}

// RedisConfig holds Redis cache configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// PackServiceConfig holds pack service gRPC connection configuration
type PackServiceConfig struct {
	Host string
	Port string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")

	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; using environment variables and defaults
	}

	cfg := &Config{
		Server: ServerConfig{
			HTTPPort: viper.GetString("HTTP_PORT"),
			WSPort:   viper.GetString("WS_PORT"),
			GRPCPort: viper.GetString("GRPC_PORT"),
			Mode:     viper.GetString("GIN_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("POSTGRES_DB"),
			SSLMode:  viper.GetString("POSTGRES_SSLMODE"),
			MaxConns: viper.GetInt("POSTGRES_MAX_CONNS"),
			MaxIdle:  viper.GetInt("POSTGRES_MAX_IDLE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		PackService: PackServiceConfig{
			Host: viper.GetString("PACK_SERVICE_HOST"),
			Port: viper.GetString("PACK_SERVICE_PORT"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func setDefaults() {
	// Server
	viper.SetDefault("HTTP_PORT", "8003")
	viper.SetDefault("WS_PORT", "8083")
	viper.SetDefault("GRPC_PORT", "50053")
	viper.SetDefault("GIN_MODE", "debug")

	// Database
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", "5435")
	viper.SetDefault("POSTGRES_USER", "gameuser")
	viper.SetDefault("POSTGRES_PASSWORD", "gamepass")
	viper.SetDefault("POSTGRES_DB", "game_db")
	viper.SetDefault("POSTGRES_SSLMODE", "disable")
	viper.SetDefault("POSTGRES_MAX_CONNS", 25)
	viper.SetDefault("POSTGRES_MAX_IDLE", 5)

	// Redis
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 2)

	// Pack Service (HTTP port, will migrate to gRPC later)
	viper.SetDefault("PACK_SERVICE_HOST", "localhost")
	viper.SetDefault("PACK_SERVICE_PORT", "8084")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("POSTGRES_HOST is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("POSTGRES_USER is required")
	}

	if c.Database.Password == "" {
		return fmt.Errorf("POSTGRES_PASSWORD is required")
	}

	if c.Redis.Host == "" {
		return fmt.Errorf("REDIS_HOST is required")
	}

	if c.PackService.Host == "" {
		return fmt.Errorf("PACK_SERVICE_HOST is required")
	}

	return nil
}

// GetPostgresConnectionString returns the PostgreSQL connection string
func (c *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetRedisAddress returns the Redis address
func (c *Config) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// GetPackServiceAddress returns the Pack Service gRPC address
func (c *Config) GetPackServiceAddress() string {
	return fmt.Sprintf("%s:%s", c.PackService.Host, c.PackService.Port)
}

// GetCacheTTL returns cache TTL for different types
func GetCacheTTL(cacheType string) time.Duration {
	switch cacheType {
	case "pack":
		return 1 * time.Hour
	case "game_state":
		return 24 * time.Hour
	default:
		return 1 * time.Hour
	}
}

