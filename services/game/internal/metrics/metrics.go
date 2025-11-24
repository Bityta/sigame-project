package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the game service
type Metrics struct {
	// HTTP metrics
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	
	// gRPC metrics
	grpcRequestsTotal   *prometheus.CounterVec
	grpcRequestDuration *prometheus.HistogramVec
	
	// WebSocket metrics
	wsConnections prometheus.Gauge
	
	// Business metrics
	activeSessions     prometheus.Gauge
	totalSessions      prometheus.Counter
	questionsAnswered  prometheus.Counter
	buttonPressLatency prometheus.Histogram
}

// NewMetrics creates and registers all metrics
func NewMetrics() *Metrics {
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
		
		// WebSocket metrics
		wsConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "game_ws_connections",
			Help: "Number of active WebSocket connections",
		}),
		
		// Business metrics
		activeSessions: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "game_sessions_active",
			Help: "Number of currently active game sessions",
		}),
		totalSessions: promauto.NewCounter(prometheus.CounterOpts{
			Name: "game_sessions_total",
			Help: "Total number of game sessions created",
		}),
		questionsAnswered: promauto.NewCounter(prometheus.CounterOpts{
			Name: "game_questions_answered_total",
			Help: "Total number of questions answered",
		}),
		buttonPressLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "game_button_press_latency_seconds",
			Help:    "Button press latency from question shown to button pressed",
			Buckets: []float64{.01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		}),
	}
}

// RecordHTTPRequest records an HTTP request
func (m *Metrics) RecordHTTPRequest(method, endpoint string, status int, duration float64) {
	statusCode := statusCodeString(status)
	m.httpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	m.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// RecordGRPCRequest records gRPC request metrics
func (m *Metrics) RecordGRPCRequest(method string, status string, duration float64) {
	m.grpcRequestsTotal.WithLabelValues(method, status).Inc()
	m.grpcRequestDuration.WithLabelValues(method).Observe(duration)
}

// RecordGameCreated records a new game creation
func (m *Metrics) RecordGameCreated() {
	m.totalSessions.Inc()
	m.activeSessions.Inc()
}

// RecordGameFinished records a game finish
func (m *Metrics) RecordGameFinished() {
	m.activeSessions.Dec()
}

// RecordWSConnection records a WebSocket connection change
func (m *Metrics) RecordWSConnection(delta int) {
	if delta > 0 {
		m.wsConnections.Inc()
	} else {
		m.wsConnections.Dec()
	}
}

// RecordButtonPress records button press latency
func (m *Metrics) RecordButtonPress(latencySeconds float64) {
	m.buttonPressLatency.Observe(latencySeconds)
}

// RecordQuestionAnswered records a question being answered
func (m *Metrics) RecordQuestionAnswered() {
	m.questionsAnswered.Inc()
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
