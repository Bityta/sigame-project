package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the game service
type Metrics struct {
	// Game sessions
	ActiveSessions prometheus.Gauge
	TotalSessions  prometheus.Counter

	// WebSocket connections
	WSConnections prometheus.Gauge

	// Button press latency
	ButtonPressLatency prometheus.Histogram

	// Questions
	QuestionsAnswered prometheus.Counter

	// Errors
	ErrorsTotal prometheus.Counter

	// HTTP requests
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
}

// NewMetrics creates and registers all metrics
func NewMetrics() *Metrics {
	return &Metrics{
		ActiveSessions: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "game_sessions_active",
			Help: "Number of currently active game sessions",
		}),

		TotalSessions: promauto.NewCounter(prometheus.CounterOpts{
			Name: "game_sessions_total",
			Help: "Total number of game sessions created",
		}),

		WSConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "game_ws_connections",
			Help: "Number of active WebSocket connections",
		}),

		ButtonPressLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "game_button_press_latency_seconds",
			Help:    "Button press latency from question shown to button pressed",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10), // 1ms to ~1s
		}),

		QuestionsAnswered: promauto.NewCounter(prometheus.CounterOpts{
			Name: "game_questions_answered_total",
			Help: "Total number of questions answered",
		}),

		ErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "game_errors_total",
			Help: "Total number of errors",
		}),

		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "game_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),

		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "game_http_request_duration_seconds",
				Help:    "HTTP request latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
	}
}

// RecordGameCreated records a new game creation
func (m *Metrics) RecordGameCreated() {
	m.TotalSessions.Inc()
	m.ActiveSessions.Inc()
}

// RecordGameFinished records a game finish
func (m *Metrics) RecordGameFinished() {
	m.ActiveSessions.Dec()
}

// RecordWSConnection records a WebSocket connection
func (m *Metrics) RecordWSConnection(delta int) {
	if delta > 0 {
		m.WSConnections.Inc()
	} else {
		m.WSConnections.Dec()
	}
}

// RecordButtonPress records button press latency
func (m *Metrics) RecordButtonPress(latencySeconds float64) {
	m.ButtonPressLatency.Observe(latencySeconds)
}

// RecordQuestionAnswered records a question being answered
func (m *Metrics) RecordQuestionAnswered() {
	m.QuestionsAnswered.Inc()
}

// RecordError records an error
func (m *Metrics) RecordError() {
	m.ErrorsTotal.Inc()
}

// RecordHTTPRequest records an HTTP request
func (m *Metrics) RecordHTTPRequest(method, endpoint, status string, duration float64) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

