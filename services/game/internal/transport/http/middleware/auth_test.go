package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authClient "sigame/game/internal/adapter/grpc/auth"
)

type MockAuthClient struct {
	mock.Mock
}

func (m *MockAuthClient) ValidateToken(ctx context.Context, token string) (*authClient.ValidateTokenResponse, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authClient.ValidateTokenResponse), args.Error(1)
}

func TestAuthMiddleware_WithValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	mockClient := new(MockAuthClient)
	mockClient.On("ValidateToken", mock.Anything, "valid-token").Return(&authClient.ValidateTokenResponse{
		Valid:  true,
		UserID: userID,
	}, nil)

	SetAuthClient(mockClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/game/my-active", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")

	handler := Auth()
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	val, exists := c.Get(UserIDContextKey)
	assert.True(t, exists)
	assert.Equal(t, userID, val)
	mockClient.AssertExpectations(t)
}

func TestAuthMiddleware_WithInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := new(MockAuthClient)
	mockClient.On("ValidateToken", mock.Anything, "invalid-token").Return(&authClient.ValidateTokenResponse{
		Valid: false,
		Error: "token expired",
	}, nil)

	SetAuthClient(mockClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/game/my-active", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	handler := Auth()
	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockClient.AssertExpectations(t)
}

func TestAuthMiddleware_WithXUserIDHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	mockClient := new(MockAuthClient)

	SetAuthClient(mockClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/game", nil)
	c.Request.Header.Set("X-User-ID", userID.String())

	handler := Auth()
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	val, exists := c.Get(UserIDContextKey)
	assert.True(t, exists)
	assert.Equal(t, userID, val)
}

func TestAuthMiddleware_NoAuthHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := new(MockAuthClient)
	SetAuthClient(mockClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/game/my-active", nil)

	handler := Auth()
	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidXUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := new(MockAuthClient)
	SetAuthClient(mockClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/game", nil)
	c.Request.Header.Set("X-User-ID", "invalid-uuid")

	handler := Auth()
	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

