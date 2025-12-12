package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	grpcRequestsTotal   *prometheus.CounterVec
	grpcRequestDuration *prometheus.HistogramVec
	wsConnections       prometheus.Gauge
	activeSessions      prometheus.Gauge
	totalSessions       prometheus.Counter
	questionsAnswered   prometheus.Counter
	buttonPressLatency  prometheus.Histogram
}

func NewMetrics() *Metrics {
	httpTotal, httpDuration := buildHTTPMetrics()
	grpcTotal, grpcDuration := buildGRPCMetrics()
	wsConn := buildWebSocketMetrics()
	activeSessions, totalSessions, questionsAnswered, buttonLatency := buildBusinessMetrics()

	return &Metrics{
		httpRequestsTotal:   httpTotal,
		httpRequestDuration: httpDuration,
		grpcRequestsTotal:   grpcTotal,
		grpcRequestDuration: grpcDuration,
		wsConnections:       wsConn,
		activeSessions:      activeSessions,
		totalSessions:       totalSessions,
		questionsAnswered:   questionsAnswered,
		buttonPressLatency:  buttonLatency,
	}
}

func (m *Metrics) RecordHTTPRequest(method, endpoint string, status int, duration float64) {
	statusCode := statusCodeString(status)
	m.httpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	m.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

func (m *Metrics) RecordGRPCRequest(method string, status string, duration float64) {
	m.grpcRequestsTotal.WithLabelValues(method, status).Inc()
	m.grpcRequestDuration.WithLabelValues(method).Observe(duration)
}

func (m *Metrics) RecordGameCreated() {
	m.totalSessions.Inc()
	m.activeSessions.Inc()
}

func (m *Metrics) RecordGameFinished() {
	m.activeSessions.Dec()
}

func (m *Metrics) RecordWSConnection(delta int) {
	if delta > 0 {
		m.wsConnections.Inc()
	} else {
		m.wsConnections.Dec()
	}
}

func (m *Metrics) RecordButtonPress(latencySeconds float64) {
	m.buttonPressLatency.Observe(latencySeconds)
}

func (m *Metrics) RecordQuestionAnswered() {
	m.questionsAnswered.Inc()
}
