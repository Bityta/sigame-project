package logger

import (
	"context"
	"testing"
)

func TestInit(t *testing.T) {
	Init("test-service")
	if serviceName != "test-service" {
		t.Errorf("Init() serviceName = %s, want test-service", serviceName)
	}
}

func TestLogFunctions(t *testing.T) {
	Init("test-service")
	ctx := context.Background()

	Info(ctx, "test info message")
	Debug(ctx, "test debug message")
	Warn(ctx, "test warn message")
	Error(ctx, "test error message")
}

func TestLogfFunctions(t *testing.T) {
	Init("test-service")
	ctx := context.Background()

	Infof(ctx, "test info message: %s", "value")
	Debugf(ctx, "test debug message: %d", 123)
	Warnf(ctx, "test warn message: %v", true)
	Errorf(ctx, "test error message: %s", "error")
}

func TestLogEntryf_Formatting(t *testing.T) {
	Init("test-service")
	ctx := context.Background()

	logEntryf(ctx, LevelInfo, "formatted: %s %d", "test", 42)
	logEntryf(ctx, LevelInfo, "no args")
}

