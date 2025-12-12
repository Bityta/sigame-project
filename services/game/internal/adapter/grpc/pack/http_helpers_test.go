package pack

import (
	"testing"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		args     []interface{}
		expected string
	}{
		{
			name:     "pack content URL",
			baseURL:  "http:
			path:     PathPackContent,
			args:     []interface{}{"123e4567-e89b-12d3-a456-426614174000"},
			expected: "http:
		},
		{
			name:     "pack URL",
			baseURL:  "http:
			path:     PathPack,
			args:     []interface{}{"123e4567-e89b-12d3-a456-426614174000"},
			expected: "http:
		},
		{
			name:     "no args",
			baseURL:  "http:
			path:     "/api/test",
			args:     []interface{}{},
			expected: "http:
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildURL(tt.baseURL, tt.path, tt.args...)
			if got != tt.expected {
				t.Errorf("buildURL() = %s, want %s", got, tt.expected)
			}
		})
	}
}

