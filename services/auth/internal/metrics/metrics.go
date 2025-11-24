package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds Prometheus metrics for the auth service
type Metrics struct {
	// HTTP metrics
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	
	// gRPC metrics
	grpcRequestsTotal   *prometheus.CounterVec
	grpcRequestDuration *prometheus.HistogramVec
	
	// Business metrics
	activeSessions prometheus.Gauge
	totalUsers     prometheus.Gauge
}

// New creates a new Metrics instance
func New() *Metrics {
	return &Metrics{
		// HTTP metrics
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency in seconds",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "endpoint"},
		),
		
		// gRPC metrics
		grpcRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_requests_total",
				Help: "Total number of gRPC requests",
			},
			[]string{"method", "status"},
		),
		grpcRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_request_duration_seconds",
				Help:    "gRPC request latency in seconds",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method"},
		),
		
		// Business metrics
		activeSessions: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "auth_active_sessions",
				Help: "Current number of active user sessions",
			},
		),
		totalUsers: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "auth_total_users",
				Help: "Total number of registered users",
			},
		),
	}
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(method, path string, status int, duration float64) {
	statusCode := statusCodeString(status)
	m.httpRequestsTotal.WithLabelValues(method, path, statusCode).Inc()
	m.httpRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// RecordGRPCRequest records gRPC request metrics
func (m *Metrics) RecordGRPCRequest(method string, status string, duration float64) {
	m.grpcRequestsTotal.WithLabelValues(method, status).Inc()
	m.grpcRequestDuration.WithLabelValues(method).Observe(duration)
}

// SetActiveSessions sets active sessions count
func (m *Metrics) SetActiveSessions(count int) {
	m.activeSessions.Set(float64(count))
}

// IncActiveSessions increments active sessions
func (m *Metrics) IncActiveSessions() {
	m.activeSessions.Inc()
}

// DecActiveSessions decrements active sessions
func (m *Metrics) DecActiveSessions() {
	m.activeSessions.Dec()
}

// SetTotalUsers sets total users count
func (m *Metrics) SetTotalUsers(count int) {
	m.totalUsers.Set(float64(count))
}

// IncTotalUsers increments total users
func (m *Metrics) IncTotalUsers() {
	m.totalUsers.Inc()
}

// statusCodeString returns status code as string (2xx/3xx/4xx/5xx)
func statusCodeString(status int) string {
	if status >= 200 && status < 300 {
		return "2xx"
	} else if status >= 300 && status < 400 {
		return "3xx"
	} else if status >= 400 && status < 500 {
		return "4xx"
	} else if status >= 500 {
		return "5xx"
	}
	return "unknown"
}

