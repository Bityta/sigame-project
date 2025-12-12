package handler

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
	"sigame/game/internal/transport/ws/hub"
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

type MockGameManager struct {
	mock.Mock
}

func (m *MockGameManager) HandleClientMessage(userID uuid.UUID, message interface{}) {
	m.Called(userID, message)
}

func (m *MockGameManager) SendStateToClient(client interface{}) {
	m.Called(client)
}

func (m *MockGameManager) SetPlayerConnected(userID uuid.UUID, connected bool) {
	m.Called(userID, connected)
}

func (m *MockGameManager) Stop() {
	m.Called()
}

func TestHandler_HandleWebSocket_WithValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gameID := uuid.New()
	userID := uuid.New()

	mockHub := hub.New()
	mockManager := new(MockGameManager)
	mockHub.RegisterGameManager(gameID, mockManager)

	mockAuthClient := new(MockAuthClient)
	mockAuthClient.On("ValidateToken", mock.Anything, "valid-token").Return(&authClient.ValidateTokenResponse{
		Valid:  true,
		UserID: userID,
	}, nil)

	h := NewHandler(mockHub, mockAuthClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ws/game/"+gameID.String()+"?token=valid-token", nil)
	c.Params = gin.Params{{Key: "id", Value: gameID.String()}}

	h.HandleWebSocket(c)

	assert.NotEqual(t, http.StatusUnauthorized, w.Code)
	mockAuthClient.AssertExpectations(t)
}

func TestHandler_HandleWebSocket_WithInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gameID := uuid.New()

	mockHub := hub.New()
	mockManager := new(MockGameManager)
	mockHub.RegisterGameManager(gameID, mockManager)

	mockAuthClient := new(MockAuthClient)
	mockAuthClient.On("ValidateToken", mock.Anything, "invalid-token").Return(&authClient.ValidateTokenResponse{
		Valid: false,
		Error: "token expired",
	}, nil)

	h := NewHandler(mockHub, mockAuthClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ws/game/"+gameID.String()+"?token=invalid-token", nil)
	c.Params = gin.Params{{Key: "id", Value: gameID.String()}}

	h.HandleWebSocket(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuthClient.AssertExpectations(t)
}

func TestHandler_HandleWebSocket_WithUserIDFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gameID := uuid.New()
	userID := uuid.New()

	mockHub := hub.New()
	mockManager := new(MockGameManager)
	mockHub.RegisterGameManager(gameID, mockManager)

	h := NewHandler(mockHub, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ws/game/"+gameID.String()+"?user_id="+userID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: gameID.String()}}

	h.HandleWebSocket(c)

	assert.NotEqual(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_HandleWebSocket_InvalidGameID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockHub := hub.New()
	mockManager := new(MockGameManager)
	mockHub.RegisterGameManager(uuid.New(), mockManager)
	h := NewHandler(mockHub, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ws/game/invalid-id", nil)
	c.Params = gin.Params{{Key: "id", Value: "invalid-id"}}

	h.HandleWebSocket(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_HandleWebSocket_GameNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gameID := uuid.New()
	userID := uuid.New()

	mockHub := hub.New()

	mockAuthClient := new(MockAuthClient)
	mockAuthClient.On("ValidateToken", mock.Anything, "valid-token").Return(&authClient.ValidateTokenResponse{
		Valid:  true,
		UserID: userID,
	}, nil)

	h := NewHandler(mockHub, mockAuthClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ws/game/"+gameID.String()+"?token=valid-token", nil)
	c.Params = gin.Params{{Key: "id", Value: gameID.String()}}

	h.HandleWebSocket(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockAuthClient.AssertExpectations(t)
}

func TestHandler_HandleWebSocket_NoTokenOrUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gameID := uuid.New()

	mockHub := hub.New()
	mockManager := new(MockGameManager)
	mockHub.RegisterGameManager(gameID, mockManager)

	h := NewHandler(mockHub, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ws/game/"+gameID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: gameID.String()}}

	h.HandleWebSocket(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

