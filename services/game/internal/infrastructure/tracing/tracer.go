package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"sigame/game/internal/infrastructure/logger"
)

func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	exporter, err := connectWithRetry()
	if err != nil || exporter == nil {
		return nil, nil
	}

	res, err := createResource(serviceName)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logger.Infof(nil, "OpenTelemetry tracer initialized for %s", serviceName)
	return tp, nil
}

func Shutdown(tp *sdktrace.TracerProvider) {
	if tp != nil {
		ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			logger.Errorf(nil, "Error shutting down tracer provider: %v", err)
		}
	}
}

type TracerProvider = sdktrace.TracerProvider


