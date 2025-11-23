package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	
	"github.com/sigame/auth/internal/domain"
)

type Metrics struct {
	requestsTotal          *prometheus.CounterVec
	requestDuration        *prometheus.HistogramVec
	jwtGenerationDuration  prometheus.Histogram
	usernameChecksTotal    *prometheus.CounterVec
	loginAttemptsTotal     *prometheus.CounterVec
	activeSessions         prometheus.Gauge
	totalUsers             prometheus.Gauge
	grpcRequestsTotal      *prometheus.CounterVec
	grpcRequestDuration    *prometheus.HistogramVec
}

func New() *Metrics {
	return &Metrics{
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "auth_requests_total",
				Help: "Total number of auth service requests",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "auth_request_duration_seconds",
				Help:    "Duration of auth service requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		jwtGenerationDuration: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "auth_jwt_generation_duration_seconds",
				Help:    "Duration of JWT token generation in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1},
			},
		),
		usernameChecksTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "auth_username_checks_total",
				Help: "Total number of username availability checks",
			},
			[]string{"available"},
		),
		loginAttemptsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "auth_login_attempts_total",
				Help: "Total number of login attempts",
			},
			[]string{"success"},
		),
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
		grpcRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "auth_grpc_requests_total",
				Help: "Total number of gRPC requests",
			},
			[]string{"method", "status"},
		),
		grpcRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "auth_grpc_request_duration_seconds",
				Help:    "Duration of gRPC requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method"},
		),
	}
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(method, path string, status int, duration float64) {
	m.requestsTotal.WithLabelValues(method, path, statusToString(status)).Inc()
	m.requestDuration.WithLabelValues(method, path).Observe(duration)
}

// RecordJWTGeneration records JWT token generation duration
func (m *Metrics) RecordJWTGeneration(duration float64) {
	m.jwtGenerationDuration.Observe(duration)
}

// RecordUsernameCheck records username availability check
func (m *Metrics) RecordUsernameCheck(available bool) {
	m.usernameChecksTotal.WithLabelValues(boolToString(available)).Inc()
}

// RecordLoginAttempt records login attempt
func (m *Metrics) RecordLoginAttempt(success bool) {
	m.loginAttemptsTotal.WithLabelValues(boolToString(success)).Inc()
}

// SetActiveSessions sets the current number of active sessions
func (m *Metrics) SetActiveSessions(count int) {
	m.activeSessions.Set(float64(count))
}

// IncActiveSessions increments active sessions counter
func (m *Metrics) IncActiveSessions() {
	m.activeSessions.Inc()
}

// DecActiveSessions decrements active sessions counter
func (m *Metrics) DecActiveSessions() {
	m.activeSessions.Dec()
}

// SetTotalUsers sets the total number of registered users
func (m *Metrics) SetTotalUsers(count int) {
	m.totalUsers.Set(float64(count))
}

// IncTotalUsers increments total users counter
func (m *Metrics) IncTotalUsers() {
	m.totalUsers.Inc()
}

// RecordGRPCRequest records gRPC request metrics
func (m *Metrics) RecordGRPCRequest(method string, status string, duration float64) {
	m.grpcRequestsTotal.WithLabelValues(method, status).Inc()
	m.grpcRequestDuration.WithLabelValues(method).Observe(duration)
}

func statusToString(status int) string {
	if status >= 200 && status < 300 {
		return domain.StatusCode2xx
	} else if status >= 300 && status < 400 {
		return domain.StatusCode3xx
	} else if status >= 400 && status < 500 {
		return domain.StatusCode4xx
	} else if status >= 500 {
		return domain.StatusCode5xx
	}
	return domain.StatusCodeUnknown
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

