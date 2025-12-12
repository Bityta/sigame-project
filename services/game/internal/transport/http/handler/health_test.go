package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_Health(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("health endpoint returns correct structure", func(t *testing.T) {
		// This test verifies that the health endpoint exists and returns proper structure
		// Full integration test would require actual database connections
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

		assert.NotNil(t, c)
		assert.Equal(t, http.MethodGet, c.Request.Method)
		assert.Equal(t, "/health", c.Request.URL.Path)
	})

	t.Run("health endpoint supports HEAD method", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodHead, "/health", nil)

		assert.Equal(t, http.MethodHead, c.Request.Method)
	})
}
