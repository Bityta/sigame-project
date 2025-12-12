package config

import "fmt"

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

func (c *Config) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

func (c *Config) GetPackServiceAddress() string {
	return fmt.Sprintf("%s:%s", c.PackService.Host, c.PackService.Port)
}

func (c *Config) GetAuthServiceAddress() string {
	return fmt.Sprintf("%s:%s", c.AuthService.Host, c.AuthService.Port)
}

