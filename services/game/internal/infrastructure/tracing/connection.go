package tracing

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/sigame/game/internal/logger"
)

func connectWithRetry() (*otlptrace.Exporter, error) {
	endpoint := os.Getenv(EnvOTLPEndpoint)
	if endpoint == "" {
		endpoint = DefaultOTLPEndpoint
	}

	var exporter *otlptrace.Exporter
	var err error

	for i := 0; i < MaxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)

		conn, connErr := grpc.DialContext(
			ctx,
			endpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)

		if connErr == nil {
			exporter, err = createExporter(ctx, conn)
			cancel()
			if err == nil {
				logger.Infof(nil, "Connected to Tempo at %s (attempt %d/%d)", endpoint, i+1, MaxRetries)
				return exporter, nil
			}
		}

		cancel()

		if i < MaxRetries-1 {
			waitTime := RetryBackoffBase * time.Duration(i+1)
			logger.Warnf(nil, "Failed to connect to Tempo (attempt %d/%d), retrying in %v...", i+1, MaxRetries, waitTime)
			time.Sleep(waitTime)
		} else {
			logger.Warnf(nil, "Could not connect to Tempo after %d attempts, continuing without tracing", MaxRetries)
			return nil, nil
		}
	}

	return nil, nil
}

func createExporter(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

