package config

import "time"

type CacheType string

const (
	CacheTypePack     CacheType = "pack"
	CacheTypeGameState CacheType = "game_state"
)

const (
	PackCacheTTL      = 1 * time.Hour
	GameStateCacheTTL = 24 * time.Hour
)

const (
	keyHTTPPort = "HTTP_PORT"
	keyWSPort   = "WS_PORT"
	keyGRPCPort = "GRPC_PORT"

	keyPostgresHost     = "POSTGRES_HOST"
	keyPostgresPort     = "POSTGRES_PORT"
	keyPostgresUser     = "POSTGRES_USER"
	keyPostgresPassword = "POSTGRES_PASSWORD"
	keyPostgresDB       = "POSTGRES_DB"
	keyPostgresSSLMode  = "POSTGRES_SSLMODE"
	keyPostgresMaxConns = "POSTGRES_MAX_CONNS"
	keyPostgresMaxIdle  = "POSTGRES_MAX_IDLE"

	keyRedisHost     = "REDIS_HOST"
	keyRedisPort     = "REDIS_PORT"
	keyRedisPassword = "REDIS_PASSWORD"
	keyRedisDB       = "REDIS_DB"

	keyPackServiceHost = "PACK_SERVICE_HOST"
	keyPackServicePort = "PACK_SERVICE_PORT"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	PackService PackServiceConfig
}

type ServerConfig struct {
	HTTPPort string
	WSPort   string
	GRPCPort string
}

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

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type PackServiceConfig struct {
	Host string
	Port string
}

