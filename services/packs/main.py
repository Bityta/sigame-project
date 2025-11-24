"""Pack Service - Main application entry point"""
import asyncio
import signal
import sys
import logging
import threading
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicorn

from app.config import settings
from app.api.routes import router
from app.metrics import metrics
from app.mock_data import MOCK_PACKS
from app.grpc.server import serve_grpc
from app.tracing import init_tracer, instrument_fastapi, shutdown_tracer
from app.middleware import AsyncRequestResponseLoggingMiddleware

# Configure logging with JSON format for better parsing
logging.basicConfig(
    level=logging.INFO,
    format='{"timestamp":"%(asctime)s","level":"%(levelname)s","service":"pack-service","message":"%(message)s"}'
)
logger = logging.getLogger(__name__)

# Global gRPC server reference
grpc_server = None
tracer_provider = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifecycle manager for FastAPI app"""
    global tracer_provider
    
    # Startup
    logger.info("Starting Pack Service...")
    logger.info(f"Service: {settings.service_name} v{settings.version}")
    logger.info(f"HTTP Server: :{settings.http_port}")
    logger.info(f"gRPC Server: :{settings.grpc_port}")
    
    # Tracer is already initialized before app creation
    # Just reference the global variable
    logger.info("✓ OpenTelemetry tracer already initialized")
    
    # Initialize metrics
    metrics.set_available_packs(len(MOCK_PACKS))
    logger.info(f"Loaded {len(MOCK_PACKS)} mock packs")
    
    # Start gRPC server in a separate thread
    global grpc_server
    grpc_thread = threading.Thread(target=start_grpc_server, daemon=True)
    grpc_thread.start()
    
    logger.info("✓ Pack Service started successfully")
    logger.info("✓ Ready to accept requests")
    
    yield
    
    # Shutdown
    logger.info("Shutting down Pack Service...")
    
    # Shutdown tracer
    shutdown_tracer(tracer_provider)
    
    if grpc_server:
        grpc_server.stop(grace=5)
        logger.info("✓ gRPC server stopped")
    
    logger.info("✓ Pack Service stopped")


def start_grpc_server():
    """Start gRPC server in a separate thread"""
    global grpc_server
    try:
        grpc_server = serve_grpc(settings.grpc_port)
        if grpc_server:
            grpc_server.wait_for_termination()
    except Exception as e:
        logger.error(f"Failed to start gRPC server: {e}")


# Create FastAPI app
app = FastAPI(
    title="Pack Service",
    description="Mock Pack Management Service for SIGame 2.0",
    version=settings.version,
    lifespan=lifespan
)

# Initialize tracer early (before adding middleware)
# We need to initialize it here to instrument FastAPI properly
try:
    tracer_provider = init_tracer("pack-service")
    # Instrument FastAPI BEFORE adding CORS middleware
    instrument_fastapi(app)
    logger.info("✓ FastAPI instrumented with OpenTelemetry")
except Exception as e:
    logger.warning(f"⚠️  Failed to instrument FastAPI: {e}")

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # In production, specify actual origins
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Add async request/response logging middleware
app.add_middleware(AsyncRequestResponseLoggingMiddleware)

# Include routes
app.include_router(router)


def signal_handler(signum, frame):
    """Handle shutdown signals"""
    logger.info(f"Received signal {signum}")
    sys.exit(0)


if __name__ == "__main__":
    # Register signal handlers
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)
    
    # Start HTTP server
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=settings.http_port,
        log_level="info",
        access_log=True
    )

