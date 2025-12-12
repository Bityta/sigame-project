package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func buildHTTPMetrics() (*prometheus.CounterVec, *prometheus.HistogramVec) {
	httpRequestsTotal := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricHTTPRequestsTotal,
			Help: HelpHTTPRequestsTotal,
		},
		[]string{LabelMethod, LabelEndpoint, LabelStatus},
	)

	httpRequestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    MetricHTTPRequestDuration,
			Help:    HelpHTTPRequestDuration,
			Buckets: DefaultHTTPBuckets,
		},
		[]string{LabelMethod, LabelEndpoint},
	)

	return httpRequestsTotal, httpRequestDuration
}

func buildGRPCMetrics() (*prometheus.CounterVec, *prometheus.HistogramVec) {
	grpcRequestsTotal := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: MetricGRPCRequestsTotal,
			Help: HelpGRPCRequestsTotal,
		},
		[]string{LabelMethod, LabelStatus},
	)

	grpcRequestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    MetricGRPCRequestDuration,
			Help:    HelpGRPCRequestDuration,
			Buckets: DefaultHTTPBuckets,
		},
		[]string{LabelMethod},
	)

	return grpcRequestsTotal, grpcRequestDuration
}

func buildWebSocketMetrics() prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Name: MetricGameWSConnections,
		Help: HelpGameWSConnections,
	})
}

func buildBusinessMetrics() (prometheus.Gauge, prometheus.Counter, prometheus.Counter, prometheus.Histogram) {
	activeSessions := promauto.NewGauge(prometheus.GaugeOpts{
		Name: MetricGameSessionsActive,
		Help: HelpGameSessionsActive,
	})

	totalSessions := promauto.NewCounter(prometheus.CounterOpts{
		Name: MetricGameSessionsTotal,
		Help: HelpGameSessionsTotal,
	})

	questionsAnswered := promauto.NewCounter(prometheus.CounterOpts{
		Name: MetricGameQuestionsAnswered,
		Help: HelpGameQuestionsAnswered,
	})

	buttonPressLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    MetricGameButtonPressLatency,
		Help:    HelpGameButtonPressLatency,
		Buckets: DefaultButtonPressBuckets,
	})

	return activeSessions, totalSessions, questionsAnswered, buttonPressLatency
}

