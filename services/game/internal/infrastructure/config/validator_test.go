package config

import "testing"

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
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
					DB:   0,
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
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
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "missing WS port",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database host",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host:     "",
					User:     "testuser",
					Password: "testpass",
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database user",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "",
					Password: "testpass",
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database password",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "",
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid MaxConns",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
					MaxConns: 0,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid MaxIdle greater than MaxConns",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
					MaxConns: 10,
					MaxIdle:  15,
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "missing redis host",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "",
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid redis DB",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
					DB:   -1,
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "missing pack service host",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				PackService: PackServiceConfig{
					Host: "",
					Port: "50055",
				},
			},
			wantErr: true,
		},
		{
			name: "missing pack service port",
			config: Config{
				Server: ServerConfig{
					HTTPPort: "8003",
					WSPort:   "8083",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					User:     "testuser",
					Password: "testpass",
					MaxConns: 25,
					MaxIdle:  5,
				},
				Redis: RedisConfig{
					Host: "localhost",
				},
				PackService: PackServiceConfig{
					Host: "localhost",
					Port: "",
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

