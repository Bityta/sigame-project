package tracing

import "time"

const (
	EnvOTLPEndpoint = "OTEL_EXPORTER_OTLP_ENDPOINT"
)

const (
	DefaultOTLPEndpoint   = "tempo:4317"
	DefaultServiceVersion = "1.0.0"
)

const (
	ConnectionTimeout = 3 * time.Second
	ShutdownTimeout   = 5 * time.Second
)

const (
	MaxRetries      = 5
	RetryBackoffBase = 1 * time.Second
)

