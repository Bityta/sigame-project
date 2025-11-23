package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/sigame/auth/internal/metrics"
)

// MetricsInterceptor creates a unary server interceptor for metrics
func MetricsInterceptor(m *metrics.Metrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		// Record metrics
		duration := time.Since(start).Seconds()
		statusCode := "success"
		if err != nil {
			statusCode = status.Code(err).String()
		}

		m.RecordGRPCRequest(info.FullMethod, statusCode, duration)

		return resp, err
	}
}

