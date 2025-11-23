"""OpenTelemetry tracing configuration for Pack Service"""
import os
import logging

from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource, SERVICE_NAME, SERVICE_VERSION
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.grpc import GrpcInstrumentorServer

logger = logging.getLogger(__name__)


def init_tracer(service_name: str = "pack-service") -> TracerProvider:
    """Initialize OpenTelemetry tracer with retry logic"""
    
    # Get OTLP endpoint from environment
    endpoint = os.getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://tempo:4317")
    
    # Create resource with service information
    resource = Resource(attributes={
        SERVICE_NAME: service_name,
        SERVICE_VERSION: "1.0.0"
    })
    
    # Create tracer provider
    provider = TracerProvider(resource=resource)
    
    # Try to create OTLP exporter with retries
    max_retries = 5
    for attempt in range(max_retries):
        try:
            otlp_exporter = OTLPSpanExporter(
                endpoint=endpoint,
                insecure=True  # For local development
            )
            
            # Add batch span processor
            provider.add_span_processor(BatchSpanProcessor(otlp_exporter))
            
            # Set as global tracer provider
            trace.set_tracer_provider(provider)
            
            logger.info(f"✓ OpenTelemetry tracer initialized for {service_name} (attempt {attempt+1}/{max_retries})")
            logger.info(f"  OTLP endpoint: {endpoint}")
            break
            
        except Exception as e:
            if attempt < max_retries - 1:
                wait_time = (attempt + 1)
                logger.warning(f"⏳ Failed to initialize tracer (attempt {attempt+1}/{max_retries}), retrying in {wait_time}s: {e}")
                import time
                time.sleep(wait_time)
            else:
                logger.warning(f"⚠️  Could not connect to Tempo after {max_retries} attempts, continuing without tracing: {e}")
    
    return provider


def instrument_fastapi(app):
    """Instrument FastAPI application"""
    FastAPIInstrumentor.instrument_app(app)
    logger.info("✓ FastAPI instrumented with OpenTelemetry")


def instrument_grpc(server):
    """Instrument gRPC server"""
    GrpcInstrumentorServer().instrument_server(server)
    logger.info("✓ gRPC server instrumented with OpenTelemetry")


def get_tracer(name: str = __name__):
    """Get a tracer instance"""
    return trace.get_tracer(name)


def shutdown_tracer(provider: TracerProvider):
    """Gracefully shutdown tracer provider"""
    if provider:
        provider.shutdown()
        logger.info("✓ Tracer provider shut down")

