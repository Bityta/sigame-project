"""Prometheus metrics for Pack Service"""
from prometheus_client import Counter, Histogram, Gauge, generate_latest, CONTENT_TYPE_LATEST
from typing import Dict
import time


class Metrics:
    """Pack Service metrics"""
    
    def __init__(self):
        # HTTP request metrics
        self.http_requests_total = Counter(
            'pack_http_requests_total',
            'Total HTTP requests',
            ['method', 'endpoint', 'status']
        )
        
        self.http_request_duration_seconds = Histogram(
            'pack_http_request_duration_seconds',
            'HTTP request duration in seconds',
            ['method', 'endpoint']
        )
        
        # gRPC metrics
        self.grpc_requests_total = Counter(
            'pack_grpc_requests_total',
            'Total gRPC requests',
            ['method', 'status']
        )
        
        self.grpc_request_duration_seconds = Histogram(
            'pack_grpc_request_duration_seconds',
            'gRPC request duration in seconds',
            ['method']
        )
        
        # Business metrics
        self.pack_requests_total = Counter(
            'pack_requests_total',
            'Total pack requests by type',
            ['request_type']
        )
        
        self.available_packs = Gauge(
            'pack_available_packs_total',
            'Total number of available packs'
        )
    
    def record_http_request(self, method: str, endpoint: str, status: int, duration: float):
        """Record HTTP request metrics"""
        self.http_requests_total.labels(method=method, endpoint=endpoint, status=status).inc()
        self.http_request_duration_seconds.labels(method=method, endpoint=endpoint).observe(duration)
    
    def record_grpc_request(self, method: str, status: str, duration: float):
        """Record gRPC request metrics"""
        self.grpc_requests_total.labels(method=method, status=status).inc()
        self.grpc_request_duration_seconds.labels(method=method).observe(duration)
    
    def inc_pack_request(self, request_type: str):
        """Increment pack request counter"""
        self.pack_requests_total.labels(request_type=request_type).inc()
    
    def set_available_packs(self, count: int):
        """Set available packs gauge"""
        self.available_packs.set(count)
    
    def generate_metrics(self) -> bytes:
        """Generate Prometheus metrics"""
        return generate_latest()


# Global metrics instance
metrics = Metrics()

