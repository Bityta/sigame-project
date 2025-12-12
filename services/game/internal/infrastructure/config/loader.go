package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	setupViper()

	if err := readConfigFile(); err != nil {
		return nil, err
	}

	setDefaults()

	cfg := buildConfig()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func setupViper() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")
	viper.AutomaticEnv()
}

func readConfigFile() error {
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}
	return nil
}

func buildConfig() *Config {
	return &Config{
		Server:      buildServerConfig(),
		Database:    buildDatabaseConfig(),
		Redis:       buildRedisConfig(),
		PackService: buildPackServiceConfig(),
		AuthService: buildAuthServiceConfig(),
	}
}

func buildServerConfig() ServerConfig {
	return ServerConfig{
		HTTPPort: viper.GetString(keyHTTPPort),
		WSPort:   viper.GetString(keyWSPort),
		GRPCPort: viper.GetString(keyGRPCPort),
	}
}

func buildDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     viper.GetString(keyPostgresHost),
		Port:     viper.GetString(keyPostgresPort),
		User:     viper.GetString(keyPostgresUser),
		Password: viper.GetString(keyPostgresPassword),
		DBName:   viper.GetString(keyPostgresDB),
		SSLMode:  viper.GetString(keyPostgresSSLMode),
		MaxConns: viper.GetInt(keyPostgresMaxConns),
		MaxIdle:  viper.GetInt(keyPostgresMaxIdle),
	}
}

func buildRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     viper.GetString(keyRedisHost),
		Port:     viper.GetString(keyRedisPort),
		Password: viper.GetString(keyRedisPassword),
		DB:       viper.GetInt(keyRedisDB),
	}
}

func buildPackServiceConfig() PackServiceConfig {
	return PackServiceConfig{
		Host: viper.GetString(keyPackServiceHost),
		Port: viper.GetString(keyPackServicePort),
	}
}

func buildAuthServiceConfig() AuthServiceConfig {
	return AuthServiceConfig{
		Host: viper.GetString(keyAuthServiceHost),
		Port: viper.GetString(keyAuthServicePort),
	}
}

func setDefaults() {
	setServerDefaults()
	setDatabaseDefaults()
	setRedisDefaults()
	setPackServiceDefaults()
	setAuthServiceDefaults()
}

func setServerDefaults() {
	viper.SetDefault(keyHTTPPort, "8003")
	viper.SetDefault(keyWSPort, "8083")
	viper.SetDefault(keyGRPCPort, "50053")
}

func setDatabaseDefaults() {
	viper.SetDefault(keyPostgresHost, "localhost")
	viper.SetDefault(keyPostgresPort, "5435")
	viper.SetDefault(keyPostgresUser, "gameuser")
	viper.SetDefault(keyPostgresPassword, "gamepass")
	viper.SetDefault(keyPostgresDB, "game_db")
	viper.SetDefault(keyPostgresSSLMode, "disable")
	viper.SetDefault(keyPostgresMaxConns, 25)
	viper.SetDefault(keyPostgresMaxIdle, 5)
}

func setRedisDefaults() {
	viper.SetDefault(keyRedisHost, "localhost")
	viper.SetDefault(keyRedisPort, "6379")
	viper.SetDefault(keyRedisPassword, "")
	viper.SetDefault(keyRedisDB, 2)
}

func setPackServiceDefaults() {
	viper.SetDefault(keyPackServiceHost, "localhost")
	viper.SetDefault(keyPackServicePort, "8084")
}

func setAuthServiceDefaults() {
	viper.SetDefault(keyAuthServiceHost, "localhost")
	viper.SetDefault(keyAuthServicePort, "50051")
}

