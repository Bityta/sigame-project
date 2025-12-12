package metrics

import (
	"testing"
)

func TestStatusCodeString(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected string
	}{
		{"2xx", 200, StatusCode2xx},
		{"2xx", 201, StatusCode2xx},
		{"2xx", 299, StatusCode2xx},
		{"3xx", 300, StatusCode3xx},
		{"3xx", 301, StatusCode3xx},
		{"3xx", 399, StatusCode3xx},
		{"4xx", 400, StatusCode4xx},
		{"4xx", 404, StatusCode4xx},
		{"4xx", 499, StatusCode4xx},
		{"5xx", 500, StatusCode5xx},
		{"5xx", 503, StatusCode5xx},
		{"5xx", 599, StatusCode5xx},
		{"unknown", 100, StatusCodeUnknown},
		{"unknown", 199, StatusCodeUnknown},
		{"5xx", 600, StatusCode5xx},
		{"unknown", -1, StatusCodeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := statusCodeString(tt.status)
			if got != tt.expected {
				t.Errorf("statusCodeString(%d) = %s, want %s", tt.status, got, tt.expected)
			}
		})
	}
}

