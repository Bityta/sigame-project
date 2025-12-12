package game

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"sigame/game/internal/core/answer"
	"sigame/game/internal/core/button"
	"sigame/game/internal/core/media"
	"sigame/game/internal/core/timer"
	"sigame/game/internal/domain/event"
	domainGame "sigame/game/internal/domain/game"
	"sigame/game/internal/domain/pack"
	"sigame/game/internal/infrastructure/logger"
	"sigame/game/internal/port"
)

type Manager struct {
	game            *domainGame.Game
	pack            *pack.Pack
	hub             Hub
	ctx             context.Context
	cancel          context.CancelFunc
	actionChan      chan *PlayerAction
	timer           *timer.Timer
	timerTicker      *time.Ticker
	buttonPress     *button.Press
	mediaTracker    *media.MediaTracker
	forAllCollector *answer.ForAllCollector
	stakeInfo       *domainGame.StakeInfo
	secretTarget    *uuid.UUID
	mu              sync.RWMutex
	eventLogger     port.EventLogger
	gameRepository  port.GameRepository
	gameCache       port.GameCache
}

type PlayerAction struct {
	UserID  uuid.UUID
	Message ClientMessage
}

type ClientMessage interface {
	GetType() string
	GetPayload() map[string]interface{}
}

type Hub interface {
	Broadcast(gameID uuid.UUID, message []byte)
	GetClientRTT(gameID, userID uuid.UUID) time.Duration
}

func New(game *domainGame.Game, pack *pack.Pack, hub Hub, eventLogger port.EventLogger, gameRepository port.GameRepository, gameCache port.GameCache) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		game:            game,
		pack:            pack,
		hub:             hub,
		ctx:             ctx,
		cancel:          cancel,
		actionChan:      make(chan *PlayerAction, ManagerActionChannelBuffer),
		timer:           timer.New(),
		buttonPress:     button.New(),
		mediaTracker:    media.NewMediaTracker(InitialRoundNumber),
		forAllCollector: answer.NewForAllCollector(),
		eventLogger:     eventLogger,
		gameRepository:  gameRepository,
		gameCache:       gameCache,
	}
}

func (m *Manager) Start() {
	m.timerTicker = time.NewTicker(TimerUpdateInterval)

	go m.run()

	m.mu.Lock()
	m.startGame()
	m.mu.Unlock()
}

func (m *Manager) Stop() {
	m.cancel()
	m.timer.Stop()
	if m.timerTicker != nil {
		m.timerTicker.Stop()
	}
	m.saveGameState()
}

func (m *Manager) run() {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf(m.ctx, "Panic in game manager run loop: %v", r)
		}
	}()

	for {
		select {
		case <-m.ctx.Done():
			return

		case action := <-m.actionChan:
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Errorf(m.ctx, "Panic handling player action: %v", r)
					}
				}()
				m.handlePlayerAction(action)
			}()

		case <-m.timer.C:
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Errorf(m.ctx, "Panic handling timeout: %v", r)
					}
				}()
				m.handleTimeout()
			}()

		case <-m.timerTicker.C:
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Errorf(m.ctx, "Panic handling timer tick: %v", r)
					}
				}()
				m.handleTimerTick()
			}()
		}
	}
}

func (m *Manager) HandleClientMessage(userID uuid.UUID, msg interface{}) {
	clientMsg, ok := msg.(ClientMessage)
	if !ok {
		logger.Warnf(m.ctx, "[HandleClientMessage] Invalid message type: %T", msg)
		return
	}
	logger.Debugf(m.ctx, "[HandleClientMessage] Received: type=%s, user_id=%s, payload=%v", clientMsg.GetType(), userID, clientMsg.GetPayload())
	select {
	case m.actionChan <- &PlayerAction{UserID: userID, Message: clientMsg}:
	case <-m.ctx.Done():
	}
}

func (m *Manager) handlePlayerAction(action *PlayerAction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch action.Message.GetType() {
	case "SELECT_QUESTION":
		m.handleSelectQuestion(action)
	case "PRESS_BUTTON":
		m.handlePressButton(action.UserID)
	case "SUBMIT_ANSWER":
		m.handleSubmitAnswer(action)
	case "JUDGE_ANSWER":
		m.handleJudgeAnswer(action)
	case "MEDIA_LOAD_PROGRESS":
		m.handleMediaLoadProgress(action)
	case "MEDIA_LOAD_COMPLETE":
		m.handleMediaLoadComplete(action)
	case "TRANSFER_SECRET":
		m.handleTransferSecret(action)
	case "PLACE_STAKE":
		m.handlePlaceStake(action)
	case "SUBMIT_FOR_ALL_ANSWER":
		m.handleSubmitForAllAnswer(action)
	}
}

func (m *Manager) SetPlayerConnected(userID uuid.UUID, connected bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if player, err := m.game.GetPlayer(userID); err == nil {
		player.SetConnected(connected)
	}
}

func (m *Manager) SendStateToClient(client interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state := m.buildGameState()
	m.sendStateToClient(client, state)
}

func (m *Manager) logEvent(eventType event.Type) *event.Event {
	evt := event.New(m.game.ID, eventType)
	m.eventLogger.LogEvent(context.Background(), evt)
	return evt
}

func (m *Manager) saveGameState() {
	ctx := context.Background()

	var gameCopy *domainGame.Game
	m.mu.RLock()
	m.game.UpdatedAt = time.Now()
	gameCopy = m.game
	m.mu.RUnlock()

	go func() {
		m.saveGameStateWithRetry(ctx, gameCopy, MaxSaveRetries, SaveRetryDelay)
	}()
}

func (m *Manager) saveGameStateWithRetry(ctx context.Context, game *domainGame.Game, maxRetries int, retryDelay time.Duration) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		if err := m.gameCache.SaveGameState(ctx, game); err != nil {
			logger.Errorf(ctx, "Failed to save game state to cache (attempt %d/%d): %v", attempt+1, maxRetries, err)
			if attempt == maxRetries-1 {
				logger.Errorf(ctx, "Failed to save game state to cache after %d attempts", maxRetries)
			}
			continue
		}

		if err := m.gameRepository.UpdateGameSession(ctx, game); err != nil {
			logger.Errorf(ctx, "Failed to update game session (attempt %d/%d): %v", attempt+1, maxRetries, err)
			if attempt == maxRetries-1 {
				logger.Errorf(ctx, "Failed to update game session after %d attempts", maxRetries)
			}
			continue
		}

		return
	}
}

