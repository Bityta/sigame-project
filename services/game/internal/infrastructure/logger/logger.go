package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	Message   string `json:"message"`
}

var serviceName string

func Init(service string) {
	serviceName = service
}

func Info(ctx context.Context, message string) {
	logEntry(ctx, LevelInfo, message)
}

func Debug(ctx context.Context, message string) {
	logEntry(ctx, LevelDebug, message)
}

func Warn(ctx context.Context, message string) {
	logEntry(ctx, LevelWarn, message)
}

func Error(ctx context.Context, message string) {
	logEntry(ctx, LevelError, message)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logEntryf(ctx, LevelInfo, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logEntryf(ctx, LevelDebug, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logEntryf(ctx, LevelWarn, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logEntryf(ctx, LevelError, format, args...)
}

func logEntry(ctx context.Context, level, message string) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     level,
		Service:   serviceName,
		Message:   message,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}
	log.Println(string(jsonBytes))
}

func logEntryf(ctx context.Context, level, format string, args ...interface{}) {
	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}
	logEntry(ctx, level, message)
}

