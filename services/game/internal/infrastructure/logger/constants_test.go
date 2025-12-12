package logger

import (
	"testing"
)

func TestConstantsAreNotEmpty(t *testing.T) {
	levels := []string{
		LevelInfo,
		LevelDebug,
		LevelWarn,
		LevelError,
	}

	for _, level := range levels {
		if level == "" {
			t.Errorf("level constant is empty")
		}
	}
}

func TestConstantsAreValid(t *testing.T) {
	if LevelInfo != "INFO" {
		t.Errorf("LevelInfo = %s, want INFO", LevelInfo)
	}

	if LevelDebug != "DEBUG" {
		t.Errorf("LevelDebug = %s, want DEBUG", LevelDebug)
	}

	if LevelWarn != "WARN" {
		t.Errorf("LevelWarn = %s, want WARN", LevelWarn)
	}

	if LevelError != "ERROR" {
		t.Errorf("LevelError = %s, want ERROR", LevelError)
	}
}

