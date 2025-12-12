package game

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/core/answer"
	"github.com/sigame/game/internal/core/button"
	"github.com/sigame/game/internal/core/media"
	"github.com/sigame/game/internal/core/timer"
	"github.com/sigame/game/internal/domain/event"
	domainGame "github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/port"
)

const (
	ActionChannelBuffer   = 256
	TimerUpdateInterval   = 100 * time.Millisecond
	RoundsOverviewDuration = 10 * time.Second
	RoundIntroDuration     = 5 * time.Second
	AnswerJudgingDuration  = 30 * time.Second
	InitialRoundNumber     = 0
	RankStartIndex         = 1
)

type Manager struct {
	game            *domainGame.Game
	pack            *pack.Pack
	hub             Hub
	ctx             context.Context
	cancel          context.CancelFunc
	actionChan      chan *PlayerAction
	timer           *timer.Timer
	timerTicker     *time.Ticker
	buttonPress     *button.Press
	mediaTracker    *media.Tracker
	forAllCollector *answer.ForAllCollector
	stakeInfo       *domainGame.StakeInfo
	secretTarget    *uuid.UUID
	mu              sync.RWMutex
	eventLogger     port.EventLogger
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

func New(game *domainGame.Game, pack *pack.Pack, hub Hub, eventLogger port.EventLogger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		game:            game,
		pack:            pack,
		hub:             hub,
		ctx:             ctx,
		cancel:          cancel,
		actionChan:      make(chan *PlayerAction, ActionChannelBuffer),
		timer:           timer.New(),
		buttonPress:     button.New(),
		mediaTracker:    media.NewTracker(InitialRoundNumber),
		forAllCollector: answer.NewForAllCollector(),
		eventLogger:     eventLogger,
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
}

func (m *Manager) run() {
	for {
		select {
		case <-m.ctx.Done():
			return

		case action := <-m.actionChan:
			m.handlePlayerAction(action)

		case <-m.timer.C:
			m.handleTimeout()

		case <-m.timerTicker.C:
			m.handleTimerTick()
		}
	}
}

func (m *Manager) HandleClientMessage(userID uuid.UUID, msg ClientMessage) {
	select {
	case m.actionChan <- &PlayerAction{UserID: userID, Message: msg}:
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

