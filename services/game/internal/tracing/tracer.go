package tracing

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitTracer initializes the OpenTelemetry tracer with retry logic
func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	// Get OTLP endpoint from environment
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "tempo:4317" // default
	}

	// Try to connect with retries
	var exporter *otlptracegrpc.Exporter
	var err error
	maxRetries := 5
	
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		
		conn, connErr := grpc.DialContext(
			ctx,
			endpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		
		if connErr == nil {
			exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
			cancel()
			if err == nil {
				log.Printf("✓ Connected to Tempo at %s (attempt %d/%d)", endpoint, i+1, maxRetries)
				break
			}
		}
		
		cancel()
		
		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * time.Second
			log.Printf("⏳ Failed to connect to Tempo (attempt %d/%d), retrying in %v...", i+1, maxRetries, waitTime)
			time.Sleep(waitTime)
		} else {
			log.Printf("⚠️  Could not connect to Tempo after %d attempts, continuing without tracing", maxRetries)
			// Return nil tracer provider instead of error - service continues without tracing
			return nil, nil
		}
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator to tracecontext (W3C standard)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	log.Printf("✓ OpenTelemetry tracer initialized for %s", serviceName)
	return tp, nil
}

// Shutdown gracefully shuts down the tracer provider
func Shutdown(tp *sdktrace.TracerProvider) {
	if tp != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}

