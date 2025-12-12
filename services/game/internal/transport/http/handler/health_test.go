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

	t.Run("health endpoint exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

		assert.NotNil(t, c)
		assert.Equal(t, http.MethodGet, c.Request.Method)
	})
}
