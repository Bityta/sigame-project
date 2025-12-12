package config

import "testing"

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
		t.Errorf("GetPostgresConnectionString() = %v, want %v", connStr, expected)
	}
}

func TestConfig_GetRedisAddress(t *testing.T) {
	cfg := &Config{
		Redis: RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
	}

	addr := cfg.GetRedisAddress()
	expected := "localhost:6379"

	if addr != expected {
		t.Errorf("GetRedisAddress() = %v, want %v", addr, expected)
	}
}

func TestConfig_GetPackServiceAddress(t *testing.T) {
	cfg := &Config{
		PackService: PackServiceConfig{
			Host: "localhost",
			Port: "8084",
		},
	}

	addr := cfg.GetPackServiceAddress()
	expected := "localhost:8084"

	if addr != expected {
		t.Errorf("GetPackServiceAddress() = %v, want %v", addr, expected)
	}
}

