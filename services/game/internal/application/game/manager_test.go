package game

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	domainGame "sigame/game/internal/domain/game"
	"sigame/game/internal/domain/event"
	"sigame/game/internal/domain/pack"
	"sigame/game/internal/domain/player"
)

type MockHub struct {
	mock.Mock
}

func (m *MockHub) Broadcast(gameID uuid.UUID, message []byte) {
	m.Called(gameID, message)
}

func (m *MockHub) GetClientRTT(gameID, userID uuid.UUID) time.Duration {
	args := m.Called(gameID, userID)
	return args.Get(0).(time.Duration)
}

type MockEventLogger struct {
	mock.Mock
}

func (m *MockEventLogger) LogEvent(ctx context.Context, e *event.Event) error {
	args := m.Called(ctx, e)
	return args.Error(0)
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

type MockClientMessage struct {
	msgType  string
	payload  map[string]interface{}
}

func (m *MockClientMessage) GetType() string {
	return m.msgType
}

func (m *MockClientMessage) GetPayload() map[string]interface{} {
	return m.payload
}

func createTestGame() *domainGame.Game {
	gameID := uuid.New()
	userID := uuid.New()

	game := &domainGame.Game{
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

	p := player.New(userID, "test-player", "", player.RolePlayer)
	game.Players[userID] = p

	return game
}

func createTestPack() *pack.Pack {
	return &pack.Pack{
		ID:          uuid.New(),
		Name:        "Test Pack",
		Description: "Test Description",
		Rounds:      []*pack.Round{},
	}
}

func TestManager_Start(t *testing.T) {
	game := createTestGame()
	testPack := createTestPack()
	mockHub := new(MockHub)
	mockHub.On("Broadcast", mock.Anything, mock.Anything).Return()
	mockHub.On("GetClientRTT", mock.Anything, mock.Anything).Return(time.Duration(0))
	mockLogger := new(MockEventLogger)
	mockLogger.On("LogEvent", mock.Anything, mock.Anything).Return(nil)
	mockRepo := new(MockGameRepository)
	mockCache := new(MockGameCache)

	manager := New(game, testPack, mockHub, mockLogger, mockRepo, mockCache)

	manager.Start()

	time.Sleep(10 * time.Millisecond)

	assert.NotNil(t, manager.timerTicker)
	assert.NotNil(t, manager.ctx)

	manager.Stop()
	mockLogger.AssertExpectations(t)
	mockHub.AssertExpectations(t)
}

func TestManager_Stop(t *testing.T) {
	game := createTestGame()
	testPack := createTestPack()
	mockHub := new(MockHub)
	mockHub.On("Broadcast", mock.Anything, mock.Anything).Return()
	mockHub.On("GetClientRTT", mock.Anything, mock.Anything).Return(time.Duration(0))
	mockLogger := new(MockEventLogger)
	mockLogger.On("LogEvent", mock.Anything, mock.Anything).Return(nil)
	mockRepo := new(MockGameRepository)
	mockCache := new(MockGameCache)

	manager := New(game, testPack, mockHub, mockLogger, mockRepo, mockCache)
	manager.Start()

	time.Sleep(10 * time.Millisecond)

	manager.Stop()

	select {
	case <-manager.ctx.Done():
		assert.True(t, true)
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Context should be cancelled")
	}
	mockLogger.AssertExpectations(t)
	mockHub.AssertExpectations(t)
}

func TestManager_HandleClientMessage(t *testing.T) {
	game := createTestGame()
	testPack := createTestPack()
	mockHub := new(MockHub)
	mockHub.On("Broadcast", mock.Anything, mock.Anything).Return()
	mockHub.On("GetClientRTT", mock.Anything, mock.Anything).Return(time.Duration(0))
	mockLogger := new(MockEventLogger)
	mockLogger.On("LogEvent", mock.Anything, mock.Anything).Return(nil)
	mockRepo := new(MockGameRepository)
	mockCache := new(MockGameCache)

	manager := New(game, testPack, mockHub, mockLogger, mockRepo, mockCache)
	manager.Start()

	userID := uuid.New()
	msg := &MockClientMessage{
		msgType: "PRESS_BUTTON",
		payload: make(map[string]interface{}),
	}

	manager.HandleClientMessage(userID, msg)

	time.Sleep(50 * time.Millisecond)

	manager.Stop()
	mockLogger.AssertExpectations(t)
	mockHub.AssertExpectations(t)
}

func TestManager_SetPlayerConnected(t *testing.T) {
	game := createTestGame()
	testPack := createTestPack()
	mockHub := new(MockHub)
	mockLogger := new(MockEventLogger)
	mockRepo := new(MockGameRepository)
	mockCache := new(MockGameCache)

	manager := New(game, testPack, mockHub, mockLogger, mockRepo, mockCache)

	userID := uuid.New()
	for uid := range game.Players {
		userID = uid
		break
	}

	manager.SetPlayerConnected(userID, true)

	p, err := game.GetPlayer(userID)
	assert.NoError(t, err)
	assert.True(t, p.IsConnected)

	manager.SetPlayerConnected(userID, false)

	p, err = game.GetPlayer(userID)
	assert.NoError(t, err)
	assert.False(t, p.IsConnected)
}

func TestManager_HandleClientMessage_InvalidMessage(t *testing.T) {
	game := createTestGame()
	testPack := createTestPack()
	mockHub := new(MockHub)
	mockHub.On("Broadcast", mock.Anything, mock.Anything).Return()
	mockHub.On("GetClientRTT", mock.Anything, mock.Anything).Return(time.Duration(0))
	mockLogger := new(MockEventLogger)
	mockLogger.On("LogEvent", mock.Anything, mock.Anything).Return(nil)
	mockRepo := new(MockGameRepository)
	mockCache := new(MockGameCache)

	manager := New(game, testPack, mockHub, mockLogger, mockRepo, mockCache)
	manager.Start()

	userID := uuid.New()
	invalidMsg := "not a ClientMessage"

	manager.HandleClientMessage(userID, invalidMsg)

	time.Sleep(10 * time.Millisecond)

	manager.Stop()
	mockLogger.AssertExpectations(t)
	mockHub.AssertExpectations(t)
}

