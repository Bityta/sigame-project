package config

import "fmt"

func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}

	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("redis config: %w", err)
	}

	if err := c.PackService.Validate(); err != nil {
		return fmt.Errorf("pack service config: %w", err)
	}

	if err := c.AuthService.Validate(); err != nil {
		return fmt.Errorf("auth service config: %w", err)
	}

	return nil
}

func (s *ServerConfig) Validate() error {
	if s.HTTPPort == "" {
		return fmt.Errorf("%s is required", keyHTTPPort)
	}
	if s.WSPort == "" {
		return fmt.Errorf("%s is required", keyWSPort)
	}
	return nil
}

func (d *DatabaseConfig) Validate() error {
	if d.Host == "" {
		return fmt.Errorf("%s is required", keyPostgresHost)
	}
	if d.User == "" {
		return fmt.Errorf("%s is required", keyPostgresUser)
	}
	if d.Password == "" {
		return fmt.Errorf("%s is required", keyPostgresPassword)
	}
	if d.MaxConns <= 0 {
		return fmt.Errorf("%s must be greater than 0", keyPostgresMaxConns)
	}
	if d.MaxIdle < 0 {
		return fmt.Errorf("%s must be non-negative", keyPostgresMaxIdle)
	}
	if d.MaxIdle > d.MaxConns {
		return fmt.Errorf("%s cannot be greater than %s", keyPostgresMaxIdle, keyPostgresMaxConns)
	}
	return nil
}

func (r *RedisConfig) Validate() error {
	if r.Host == "" {
		return fmt.Errorf("%s is required", keyRedisHost)
	}
	if r.DB < 0 {
		return fmt.Errorf("%s must be non-negative", keyRedisDB)
	}
	return nil
}

func (p *PackServiceConfig) Validate() error {
	if p.Host == "" {
		return fmt.Errorf("%s is required", keyPackServiceHost)
	}
	if p.Port == "" {
		return fmt.Errorf("%s is required", keyPackServicePort)
	}
	return nil
}

func (a *AuthServiceConfig) Validate() error {
	if a.Host == "" {
		return fmt.Errorf("%s is required", keyAuthServiceHost)
	}
	if a.Port == "" {
		return fmt.Errorf("%s is required", keyAuthServicePort)
	}
	return nil
}

