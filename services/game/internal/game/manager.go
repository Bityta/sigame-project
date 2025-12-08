package game

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain"
	"github.com/sigame/game/internal/transport/websocket"
)

// Manager manages a single game session
type Manager struct {
	game          *domain.Game
	pack          *domain.Pack
	hub           Hub
	ctx           context.Context
	cancel        context.CancelFunc
	actionChan    chan *PlayerAction
	timer         *Timer
	timerTicker   *time.Ticker
	buttonPress   *ButtonPress
	mu            sync.RWMutex
	eventLogger   EventLogger
}

// PlayerAction represents an action from a player
type PlayerAction struct {
	UserID  uuid.UUID
	Message *websocket.ClientMessage
}

// Hub interface for broadcasting messages to game clients
type Hub interface {
	Broadcast(gameID uuid.UUID, message []byte)
}

// EventLogger interface for logging events
type EventLogger interface {
	LogEvent(ctx context.Context, event *domain.GameEvent) error
}

// NewManager creates a new game manager
func NewManager(game *domain.Game, pack *domain.Pack, hub Hub, eventLogger EventLogger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		game:        game,
		pack:        pack,
		hub:         hub,
		ctx:         ctx,
		cancel:      cancel,
		actionChan:  make(chan *PlayerAction, 100),
		timer:       NewTimer(),
		buttonPress: NewButtonPress(),
		eventLogger: eventLogger,
	}
}

// Start starts the game manager
func (m *Manager) Start() {
	log.Printf("Starting game manager for game %s", m.game.ID)

	// Start ticker for timer updates (every second)
	m.timerTicker = time.NewTicker(1 * time.Second)

	// Start main game loop
	go m.run()

	// Players are already ready in lobby, start game immediately
	m.mu.Lock()
	m.startGame()
	m.mu.Unlock()
}

// Stop stops the game manager
func (m *Manager) Stop() {
	log.Printf("Stopping game manager for game %s", m.game.ID)
	m.cancel()
	m.timer.Stop()
	if m.timerTicker != nil {
		m.timerTicker.Stop()
	}
}

// run is the main game loop
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
			// Broadcast timer updates during active phases
			m.handleTimerTick()
		}
	}
}

// handleTimerTick broadcasts timer updates to clients
func (m *Manager) handleTimerTick() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Only broadcast during phases that have visible timers
	switch m.game.Status {
	case domain.GameStatusQuestionSelect,
		domain.GameStatusButtonPress,
		domain.GameStatusAnswering,
		domain.GameStatusAnswerJudging:
		// Broadcast current state with updated timer
		m.BroadcastStateUnlocked()
	}
}

// BroadcastStateUnlocked broadcasts state without acquiring lock (caller must hold RLock)
func (m *Manager) BroadcastStateUnlocked() {
	state := m.buildGameState()
	msg := websocket.NewStateUpdateMessage(state)

	data, err := msg.ToJSON()
	if err != nil {
		log.Printf("Failed to serialize state: %v", err)
		return
	}

	m.hub.Broadcast(m.game.ID, data)
}

// HandleClientMessage handles a message from a client
func (m *Manager) HandleClientMessage(userID uuid.UUID, msg *websocket.ClientMessage) {
	// Put action in queue for processing
	select {
	case m.actionChan <- &PlayerAction{UserID: userID, Message: msg}:
	case <-m.ctx.Done():
	}
}

// handlePlayerAction processes a player action
func (m *Manager) handlePlayerAction(action *PlayerAction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch action.Message.Type {
	case websocket.MessageTypeSelectQuestion:
		m.handleSelectQuestion(action)

	case websocket.MessageTypePressButton:
		m.handlePressButton(action.UserID)

	case websocket.MessageTypeSubmitAnswer:
		m.handleSubmitAnswer(action)

	case websocket.MessageTypeJudgeAnswer:
		m.handleJudgeAnswer(action)

	default:
		log.Printf("Unknown message type: %s", action.Message.Type)
	}
}

// startGame starts the game
func (m *Manager) startGame() {
	log.Printf("Starting game %s", m.game.ID)

	now := time.Now()
	m.game.StartedAt = &now
	m.game.CurrentRound = 0 // Will be set to 1 when rounds_overview ends

	// Log event
	event := domain.NewGameEvent(m.game.ID, domain.EventGameStarted)
	m.eventLogger.LogEvent(context.Background(), event)

	// Show rounds overview first
	m.showRoundsOverview()
}

// showRoundsOverview displays all rounds before starting the game
func (m *Manager) showRoundsOverview() {
	log.Printf("Showing rounds overview for game %s", m.game.ID)
	m.updateGameStatus(domain.GameStatusRoundsOverview)
	m.BroadcastState()

	// Auto-transition to first round after 5 seconds
	m.timer.Start(5 * time.Second)
}

// startRound starts a new round
func (m *Manager) startRound(roundNumber int) {
	if roundNumber > len(m.pack.Rounds) {
		// All rounds complete
		m.endGame()
		return
	}

	m.game.CurrentRound = roundNumber
	m.updateGameStatus(domain.GameStatusRoundStart)

	log.Printf("Starting round %d", roundNumber)

	// Log event
	event := domain.NewGameEvent(m.game.ID, domain.EventRoundStarted).
		WithRound(roundNumber)
	m.eventLogger.LogEvent(context.Background(), event)

	// Show round intro screen
	m.BroadcastState()

	// Wait 3 seconds for round intro before showing questions
	m.timer.Start(3 * time.Second)
}

// selectActivePlayer selects the next active player (lowest score, excluding host)
func (m *Manager) selectActivePlayer() uuid.UUID {
	var selectedPlayer uuid.UUID
	minScore := int(^uint(0) >> 1) // max int

	for userID, player := range m.game.Players {
		// Skip host - they don't participate in scoring, they are the game master
		if player.Role == domain.PlayerRoleHost {
			continue
		}
		if player.IsActive && player.Score < minScore {
			minScore = player.Score
			selectedPlayer = userID
		}
	}

	return selectedPlayer
}

// findHost finds the host player
func (m *Manager) findHost() uuid.UUID {
	for userID, player := range m.game.Players {
		if player.Role == domain.PlayerRoleHost {
			return userID
		}
	}
	return uuid.Nil
}

// updateGameStatus updates the game status
func (m *Manager) updateGameStatus(status domain.GameStatus) {
	m.game.Status = status
	m.game.CurrentPhase = status
}

// BroadcastState broadcasts the current game state to all clients
func (m *Manager) BroadcastState() {
	state := m.buildGameState()
	msg := websocket.NewStateUpdateMessage(state)

	data, err := msg.ToJSON()
	if err != nil {
		log.Printf("Failed to serialize state: %v", err)
		return
	}

	m.hub.Broadcast(m.game.ID, data)
}

// SendStateToClient sends current game state to a specific client
func (m *Manager) SendStateToClient(client *websocket.Client) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state := m.buildGameState()
	msg := websocket.NewStateUpdateMessage(state)

	data, err := msg.ToJSON()
	if err != nil {
		log.Printf("Failed to serialize state for client: %v", err)
		return
	}

	client.Send(data)
	log.Printf("Sent initial state to client %s", client.GetUserID())
}

// buildGameState builds the current game state for broadcasting
func (m *Manager) buildGameState() *domain.GameState {
	state := &domain.GameState{
		GameID:        m.game.ID,
		Status:        m.game.Status,
		CurrentRound:  m.game.CurrentRound,
		Players:       make([]domain.PlayerState, 0, len(m.game.Players)),
		ActivePlayer:  m.game.ActivePlayer,
		TimeRemaining: m.timer.Remaining(),
	}

	// Add players
	for _, player := range m.game.Players {
		state.Players = append(state.Players, player.ToState())
	}

	// Add all rounds for rounds_overview status
	if m.game.Status == domain.GameStatusRoundsOverview {
		state.AllRounds = make([]domain.RoundOverview, 0, len(m.pack.Rounds))
		for i, round := range m.pack.Rounds {
			themeNames := make([]string, 0, len(round.Themes))
			for _, theme := range round.Themes {
				themeNames = append(themeNames, theme.Name)
			}
			state.AllRounds = append(state.AllRounds, domain.RoundOverview{
				RoundNumber: i + 1,
				Name:        round.Name,
				ThemeNames:  themeNames,
			})
		}
	}

	// Add round info if in game
	if m.game.CurrentRound > 0 && m.game.CurrentRound <= len(m.pack.Rounds) {
		round := m.pack.Rounds[m.game.CurrentRound-1]
		state.RoundName = round.Name

		// Add themes with question availability
		state.Themes = make([]domain.ThemeState, 0, len(round.Themes))
		for _, theme := range round.Themes {
			includeText := m.game.Status == domain.GameStatusQuestionShow
			state.Themes = append(state.Themes, theme.ToState(includeText))
		}
	}

	// Add current question if shown (include answer for host to see)
	if m.game.CurrentQuestion != nil {
		questionState := m.game.CurrentQuestion.ToStateWithAnswer(true)
		state.CurrentQuestion = &questionState
	}

	// Add winners and final scores for game_end status
	if m.game.Status == domain.GameStatusGameEnd {
		state.Winners = m.game.Winners
		state.FinalScores = m.game.FinalScores
	}

	return state
}

// handleTimeout handles timer timeout
func (m *Manager) handleTimeout() {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("Timer expired in status %s", m.game.Status)

	switch m.game.Status {
	case domain.GameStatusRoundsOverview:
		// Rounds overview finished, start first round
		m.startRound(1)

	case domain.GameStatusRoundStart:
		// Round intro finished, transition to question selection
		m.transitionToQuestionSelect()

	case domain.GameStatusQuestionSelect:
		// Auto-select random question
		m.autoSelectQuestion()

	case domain.GameStatusQuestionShow:
		// Question was shown, now allow button press
		m.transitionToButtonPress()

	case domain.GameStatusButtonPress:
		// No one pressed the button, skip question
		m.skipQuestion()

	case domain.GameStatusAnswering:
		// Time's up for answering
		m.handleAnswerTimeout()

	case domain.GameStatusAnswerJudging:
		// Host didn't judge in time - treat as wrong answer
		m.handleAnswerTimeout()
	}
}

// transitionToQuestionSelect moves from round_start to question_select
func (m *Manager) transitionToQuestionSelect() {
	log.Printf("Transitioning to question_select phase")

	// Host selects questions (they are the game master)
	hostID := m.findHost()
	m.game.ActivePlayer = &hostID

	m.updateGameStatus(domain.GameStatusQuestionSelect)

	// Start timer BEFORE broadcasting so timeRemaining is correct
	m.timer.Start(time.Duration(m.game.Settings.TimeForChoice) * time.Second)
	m.BroadcastState()
}

// transitionToButtonPress moves from question_show to button_press
func (m *Manager) transitionToButtonPress() {
	log.Printf("Transitioning to button_press phase")
	m.updateGameStatus(domain.GameStatusButtonPress)
	m.buttonPress.Reset()

	// Start timer BEFORE broadcasting so timeRemaining is correct
	m.timer.Start(time.Duration(m.game.Settings.TimeForAnswer) * time.Second)
	m.BroadcastState()
}

// autoSelectQuestion automatically selects a random available question
func (m *Manager) autoSelectQuestion() {
	round := m.pack.Rounds[m.game.CurrentRound-1]
	
	// Find available questions
	for _, theme := range round.Themes {
		for _, question := range theme.Questions {
			if question.IsAvailable() {
				// Select this question
				m.selectQuestion(theme, question)
				return
			}
		}
	}

	// No questions left in round
	m.endRound()
}

// getPlayerIDs returns a list of player IDs for debugging
func (m *Manager) getPlayerIDs() []string {
	var ids []string
	for id := range m.game.Players {
		ids = append(ids, id.String())
	}
	return ids
}

// Continue with more manager methods in the next file...

