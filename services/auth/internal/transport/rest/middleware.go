package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
	"github.com/sigame/auth/internal/domain"
	"github.com/sigame/auth/internal/metrics"
	"github.com/sigame/auth/internal/service"
)

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		token, err := ExtractBearerToken(authHeader)
		if err != nil {
			var message string
			switch err {
			case ErrAuthHeaderMissing:
				message = "Authorization header is required"
			case ErrInvalidAuthFormat:
				message = "Invalid authorization header format"
			default:
				message = "Invalid authorization"
			}
			
			respondUnauthorized(c, "unauthorized", message)
			c.Abort()
			return
		}

		// Validate token
		claims, err := authService.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			switch err {
			case domain.ErrTokenExpired:
				respondUnauthorized(c, "token_expired", "Access token has expired")
			case domain.ErrTokenBlacklisted:
				respondUnauthorized(c, "token_revoked", "Token has been revoked")
			default:
				respondUnauthorized(c, "invalid_token", "Invalid access token")
			}
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID.String())
		c.Set("username", claims.Username)

		c.Next()
	}
}

// MetricsMiddleware collects HTTP request metrics
func MetricsMiddleware(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()

		m.RecordHTTPRequest(method, path, status, duration)
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

