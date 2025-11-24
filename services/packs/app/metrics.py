"""Prometheus metrics for Pack Service"""
from prometheus_client import Counter, Histogram, Gauge, generate_latest, CONTENT_TYPE_LATEST


class Metrics:
    """Pack Service metrics"""
    
    def __init__(self):
        # HTTP metrics
        self.http_requests_total = Counter(
            'http_requests_total',
            'Total HTTP requests',
            ['method', 'endpoint', 'status']
        )
        
        self.http_request_duration_seconds = Histogram(
            'http_request_duration_seconds',
            'HTTP request duration in seconds',
            ['method', 'endpoint'],
            buckets=(.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10)
        )
        
        # gRPC metrics
        self.grpc_requests_total = Counter(
            'grpc_requests_total',
            'Total gRPC requests',
            ['method', 'status']
        )
        
        self.grpc_request_duration_seconds = Histogram(
            'grpc_request_duration_seconds',
            'gRPC request duration in seconds',
            ['method'],
            buckets=(.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10)
        )
        
        # Business metrics
        self.available_packs = Gauge(
            'packs_available_total',
            'Total number of available packs'
        )
    
    def record_http_request(self, method: str, endpoint: str, status: int, duration: float):
        """Record HTTP request metrics"""
        status_code = self._status_code_string(status)
        self.http_requests_total.labels(method=method, endpoint=endpoint, status=status_code).inc()
        self.http_request_duration_seconds.labels(method=method, endpoint=endpoint).observe(duration)
    
    def record_grpc_request(self, method: str, status: str, duration: float):
        """Record gRPC request metrics"""
        self.grpc_requests_total.labels(method=method, status=status).inc()
        self.grpc_request_duration_seconds.labels(method=method).observe(duration)
    
    def set_available_packs(self, count: int):
        """Set available packs gauge"""
        self.available_packs.set(count)
    
    def generate_metrics(self) -> bytes:
        """Generate Prometheus metrics"""
        return generate_latest()
    
    @staticmethod
    def _status_code_string(status: int) -> str:
        """Convert status code to string (2xx/3xx/4xx/5xx)"""
        if 200 <= status < 300:
            return "2xx"
        elif 300 <= status < 400:
            return "3xx"
        elif 400 <= status < 500:
            return "4xx"
        elif status >= 500:
            return "5xx"
        return "unknown"


# Global metrics instance
metrics = Metrics()
