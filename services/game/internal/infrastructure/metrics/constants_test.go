package metrics

import (
	"testing"
)

func TestMetricNamesAreNotEmpty(t *testing.T) {
	names := []string{
		MetricHTTPRequestsTotal,
		MetricHTTPRequestDuration,
		MetricGRPCRequestsTotal,
		MetricGRPCRequestDuration,
		MetricGameWSConnections,
		MetricGameSessionsActive,
		MetricGameSessionsTotal,
		MetricGameQuestionsAnswered,
		MetricGameButtonPressLatency,
	}

	for _, name := range names {
		if name == "" {
			t.Errorf("metric name constant is empty")
		}
	}
}

func TestHelpTextsAreNotEmpty(t *testing.T) {
	helps := []string{
		HelpHTTPRequestsTotal,
		HelpHTTPRequestDuration,
		HelpGRPCRequestsTotal,
		HelpGRPCRequestDuration,
		HelpGameWSConnections,
		HelpGameSessionsActive,
		HelpGameSessionsTotal,
		HelpGameQuestionsAnswered,
		HelpGameButtonPressLatency,
	}

	for _, help := range helps {
		if help == "" {
			t.Errorf("help text constant is empty")
		}
	}
}

func TestLabelsAreNotEmpty(t *testing.T) {
	labels := []string{
		LabelMethod,
		LabelEndpoint,
		LabelStatus,
	}

	for _, label := range labels {
		if label == "" {
			t.Errorf("label constant is empty")
		}
	}
}

func TestBucketsAreValid(t *testing.T) {
	if len(DefaultHTTPBuckets) == 0 {
		t.Error("DefaultHTTPBuckets should not be empty")
	}

	if len(DefaultButtonPressBuckets) == 0 {
		t.Error("DefaultButtonPressBuckets should not be empty")
	}

	for i := 1; i < len(DefaultHTTPBuckets); i++ {
		if DefaultHTTPBuckets[i] <= DefaultHTTPBuckets[i-1] {
			t.Error("DefaultHTTPBuckets should be in ascending order")
		}
	}
}

func TestStatusCodeConstants(t *testing.T) {
	codes := []string{
		StatusCode2xx,
		StatusCode3xx,
		StatusCode4xx,
		StatusCode5xx,
		StatusCodeUnknown,
	}

	for _, code := range codes {
		if code == "" {
			t.Errorf("status code constant is empty")
		}
	}
}

