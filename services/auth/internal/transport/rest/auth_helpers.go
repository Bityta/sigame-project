package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	
	"github.com/sigame/auth/internal/domain"
)

var (
	// ErrAuthHeaderMissing is returned when Authorization header is missing
	ErrAuthHeaderMissing = errors.New("authorization header is required")
	
	// ErrInvalidAuthFormat is returned when Authorization header format is invalid
	ErrInvalidAuthFormat = errors.New("invalid authorization header format, expected: Bearer <token>")
)

// ExtractBearerToken extracts the Bearer token from Authorization header
// Expected format: "Bearer <token>"
// Returns the token string or an error if the header is invalid
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrAuthHeaderMissing
	}

	// Split "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidAuthFormat
	}

	token := parts[1]
	if token == "" {
		return "", ErrInvalidAuthFormat
	}

	return token, nil
}

// respondWithError sends a JSON error response
func respondWithError(c *gin.Context, statusCode int, errorCode, message string) {
	c.JSON(statusCode, domain.ErrorResponse{
		Error:   errorCode,
		Message: message,
	})
}

// respondBadRequest sends a 400 Bad Request error response
func respondBadRequest(c *gin.Context, errorCode, message string) {
	respondWithError(c, http.StatusBadRequest, errorCode, message)
}

// respondUnauthorized sends a 401 Unauthorized error response
func respondUnauthorized(c *gin.Context, errorCode, message string) {
	respondWithError(c, http.StatusUnauthorized, errorCode, message)
}

// respondConflict sends a 409 Conflict error response
func respondConflict(c *gin.Context, errorCode, message string) {
	respondWithError(c, http.StatusConflict, errorCode, message)
}

// respondTooManyRequests sends a 429 Too Many Requests error response
func respondTooManyRequests(c *gin.Context, errorCode, message string) {
	respondWithError(c, http.StatusTooManyRequests, errorCode, message)
}

// respondInternalError sends a 500 Internal Server Error response
func respondInternalError(c *gin.Context, errorCode, message string) {
	respondWithError(c, http.StatusInternalServerError, errorCode, message)
}

// respondNotFound sends a 404 Not Found error response
func respondNotFound(c *gin.Context, errorCode, message string) {
	respondWithError(c, http.StatusNotFound, errorCode, message)
}

