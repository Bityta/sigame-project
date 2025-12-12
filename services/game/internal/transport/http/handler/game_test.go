package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	domainGame "sigame/game/internal/domain/game"
	"sigame/game/internal/domain/event"
	"sigame/game/internal/domain/pack"
	"sigame/game/internal/domain/player"
	"sigame/game/internal/transport/ws/hub"
)

type MockPackService struct {
	mock.Mock
}

func (m *MockPackService) GetPackContent(ctx context.Context, packID uuid.UUID) (*pack.Pack, error) {
	args := m.Called(ctx, packID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pack.Pack), args.Error(1)
}

func (m *MockPackService) ValidatePackExists(ctx context.Context, packID uuid.UUID) (bool, error) {
	args := m.Called(ctx, packID)
	return args.Bool(0), args.Error(1)
}

type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) CreateGameSession(ctx context.Context, g *domainGame.Game) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}

func (m *MockGameRepository) GetGameSession(ctx context.Context, gameID uuid.UUID) (*domainGame.Game, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainGame.Game), args.Error(1)
}

func (m *MockGameRepository) UpdateGameSession(ctx context.Context, g *domainGame.Game) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}

func (m *MockGameRepository) GetActiveGameForUser(ctx context.Context, userID uuid.UUID) (*domainGame.Game, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainGame.Game), args.Error(1)
}

type MockGameCache struct {
	mock.Mock
}

func (m *MockGameCache) SaveGameState(ctx context.Context, g *domainGame.Game) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}

func (m *MockGameCache) LoadGameState(ctx context.Context, gameID uuid.UUID) (*domainGame.Game, error) {
	args := m.Called(ctx, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainGame.Game), args.Error(1)
}

func (m *MockGameCache) DeleteGameState(ctx context.Context, gameID uuid.UUID) error {
	args := m.Called(ctx, gameID)
	return args.Error(0)
}

func (m *MockGameCache) GetActiveGames(ctx context.Context, limit int64) ([]uuid.UUID, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

type MockEventLogger struct {
	mock.Mock
}

func (m *MockEventLogger) LogEvent(ctx context.Context, e *event.Event) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func TestGameHandler_CreateGame(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    CreateGameRequest
		expectedStatus int
	}{
		{
			name: "missing host",
			requestBody: CreateGameRequest{
				RoomID: uuid.New(),
				PackID: uuid.New(),
				Players: []PlayerInfo{
					{UserID: uuid.New(), Username: "player", Role: "player"},
				},
				Settings: GameSettings{
					TimeForAnswer: 30,
					TimeForChoice: 20,
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid role",
			requestBody: CreateGameRequest{
				RoomID: uuid.New(),
				PackID: uuid.New(),
				Players: []PlayerInfo{
					{UserID: uuid.New(), Username: "host", Role: "invalid"},
				},
				Settings: GameSettings{
					TimeForAnswer: 30,
					TimeForChoice: 20,
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPackService := new(MockPackService)
			mockRepo := new(MockGameRepository)
			mockCache := new(MockGameCache)
			mockHub := hub.New()
			mockLogger := new(MockEventLogger)

			handler := NewGameHandler(mockPackService, mockRepo, mockCache, mockHub, mockLogger)

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateGame(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGameHandler_GetGame(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gameID := uuid.New()
	mockGame := &domainGame.Game{
		ID:           gameID,
		RoomID:       uuid.New(),
		PackID:       uuid.New(),
		Status:       domainGame.StatusWaiting,
		CurrentRound: 0,
		Players:      make(map[uuid.UUID]*player.Player),
		Settings: domainGame.Settings{
			TimeForAnswer: 30,
			TimeForChoice: 20,
		},
	}

	mockPackService := new(MockPackService)
	mockRepo := new(MockGameRepository)
	mockCache := new(MockGameCache)
	mockHub := hub.New()
	mockLogger := new(MockEventLogger)

	mockRepo.On("GetGameSession", mock.Anything, gameID).Return(mockGame, nil)

	handler := NewGameHandler(mockPackService, mockRepo, mockCache, mockHub, mockLogger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/game/"+gameID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: gameID.String()}}

	handler.GetGame(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

