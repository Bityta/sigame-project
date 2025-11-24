package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sigame/game/internal/logger"
)

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestResponseLoggingMiddleware logs request and response data at debug level (async)
func RequestResponseLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip health and metrics endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		startTime := time.Now()

		// Read request body asynchronously
		var requestBody map[string]interface{}
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// Restore body for handlers
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				
				// Try to parse JSON
				if len(bodyBytes) > 0 {
					json.Unmarshal(bodyBytes, &requestBody)
				}
			}
		}

		// Log incoming request (async via zerolog)
		logger.Log.Debug().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("client_ip", c.ClientIP()).
			Interface("body", sanitizeRequestBody(requestBody)).
			Msg("Incoming request")

		// Wrap response writer to capture response
		blw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Parse response body asynchronously
		var responseBody map[string]interface{}
		if blw.body.Len() > 0 {
			json.Unmarshal(blw.body.Bytes(), &responseBody)
		}

		// Log response (async via zerolog)
		logger.Log.Debug().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration_ms", duration).
			Interface("response", responseBody).
			Msg("Request completed")
	}
}

// sanitizeRequestBody removes sensitive fields from request body
func sanitizeRequestBody(body map[string]interface{}) map[string]interface{} {
	if body == nil {
		return nil
	}

	sanitized := make(map[string]interface{})
	for k, v := range body {
		// Keep all fields for game service (no sensitive data)
		sanitized[k] = v
	}
	return sanitized
}

