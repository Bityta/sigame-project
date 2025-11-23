package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	TraceID   string `json:"trace_id,omitempty"`
	SpanID    string `json:"span_id,omitempty"`
	Message   string `json:"message"`
}

var serviceName string

// Init initializes the logger with service name
func Init(service string) {
	serviceName = service
}

// Info logs an info message
func Info(ctx context.Context, message string) {
	logEntry(ctx, "INFO", message)
}

// Debug logs a debug message
func Debug(ctx context.Context, message string) {
	logEntry(ctx, "DEBUG", message)
}

// Warn logs a warning message
func Warn(ctx context.Context, message string) {
	logEntry(ctx, "WARN", message)
}

// Error logs an error message
func Error(ctx context.Context, message string) {
	logEntry(ctx, "ERROR", message)
}

// Infof logs an info message with formatting
func Infof(ctx context.Context, format string, args ...interface{}) {
	logEntryf(ctx, "INFO", format, args...)
}

// Debugf logs a debug message with formatting
func Debugf(ctx context.Context, format string, args ...interface{}) {
	logEntryf(ctx, "DEBUG", format, args...)
}

// Warnf logs a warning message with formatting
func Warnf(ctx context.Context, format string, args ...interface{}) {
	logEntryf(ctx, "WARN", format, args...)
}

// Errorf logs an error message with formatting
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logEntryf(ctx, "ERROR", format, args...)
}

func logEntry(ctx context.Context, level, message string) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     level,
		Service:   serviceName,
		Message:   message,
	}

	// Extract trace context if available
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		entry.TraceID = span.SpanContext().TraceID().String()
		entry.SpanID = span.SpanContext().SpanID().String()
	}

	jsonBytes, _ := json.Marshal(entry)
	log.Println(string(jsonBytes))
}

func logEntryf(ctx context.Context, level, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logEntry(ctx, level, message)
}

