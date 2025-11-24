"""
Async Request/Response Logging Middleware for Pack Service
Logs all incoming requests and outgoing responses at DEBUG level
"""

import json
import time
from typing import Callable
from fastapi import Request, Response
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.types import ASGIApp

from app.tracing import logger


class AsyncRequestResponseLoggingMiddleware(BaseHTTPMiddleware):
    """Middleware for async logging of requests and responses"""
    
    SKIP_PATHS = {"/health", "/metrics"}
    SENSITIVE_FIELDS = {"password", "token", "access_token", "refresh_token"}
    
    def __init__(self, app: ASGIApp):
        super().__init__(app)
    
    async def dispatch(self, request: Request, call_next: Callable) -> Response:
        """Process request and log async"""
        path = request.url.path
        
        # Skip health and metrics
        if path in self.SKIP_PATHS:
            return await call_next(request)
        
        start_time = time.time()
        method = request.method
        
        # Read request body asynchronously
        request_body = None
        if method in ["POST", "PUT", "PATCH"]:
            try:
                body_bytes = await request.body()
                if body_bytes:
                    request_body = json.loads(body_bytes.decode("utf-8"))
                    request_body = self._sanitize_body(request_body)
            except Exception as e:
                logger.debug(f"Failed to parse request body: {e}")
        
        # Log incoming request (async)
        logger.debug(
            "Incoming request",
            extra={
                "method": method,
                "path": path,
                "query": str(request.query_params),
                "client_ip": request.client.host if request.client else None,
                "body": request_body
            }
        )
        
        # Process request
        response = await call_next(request)
        
        # Calculate duration
        duration = (time.time() - start_time) * 1000  # ms
        
        # Log response (async)
        logger.debug(
            "Request completed",
            extra={
                "method": method,
                "path": path,
                "status": response.status_code,
                "duration_ms": f"{duration:.2f}"
            }
        )
        
        return response
    
    def _sanitize_body(self, body: dict) -> dict:
        """Remove sensitive fields from body"""
        if not isinstance(body, dict):
            return body
        
        sanitized = {}
        for key, value in body.items():
            if key.lower() in self.SENSITIVE_FIELDS:
                sanitized[key] = "***HIDDEN***"
            elif isinstance(value, dict):
                sanitized[key] = self._sanitize_body(value)
            else:
                sanitized[key] = value
        
        return sanitized

