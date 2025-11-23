package domain

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of game event
type EventType string

const (
	EventGameCreated       EventType = "GAME_CREATED"
	EventGameStarted       EventType = "GAME_STARTED"
	EventGameFinished      EventType = "GAME_FINISHED"
	EventGameCancelled     EventType = "GAME_CANCELLED"
	EventPlayerJoined      EventType = "PLAYER_JOINED"
	EventPlayerLeft        EventType = "PLAYER_LEFT"
	EventPlayerReady       EventType = "PLAYER_READY"
	EventRoundStarted      EventType = "ROUND_STARTED"
	EventRoundFinished     EventType = "ROUND_FINISHED"
	EventQuestionSelected  EventType = "QUESTION_SELECTED"
	EventQuestionShown     EventType = "QUESTION_SHOWN"
	EventButtonPressed     EventType = "BUTTON_PRESSED"
	EventAnswerSubmitted   EventType = "ANSWER_SUBMITTED"
	EventAnswerCorrect     EventType = "ANSWER_CORRECT"
	EventAnswerIncorrect   EventType = "ANSWER_INCORRECT"
	EventScoreChanged      EventType = "SCORE_CHANGED"
	EventTimerStarted      EventType = "TIMER_STARTED"
	EventTimerExpired      EventType = "TIMER_EXPIRED"
)

// GameEvent represents an event that occurred during the game
type GameEvent struct {
	ID           uuid.UUID              `json:"id"`
	GameID       uuid.UUID              `json:"game_id"`
	EventType    EventType              `json:"event_type"`
	UserID       *uuid.UUID             `json:"user_id,omitempty"`
	RoundNumber  *int                   `json:"round_number,omitempty"`
	QuestionID   *string                `json:"question_id,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Timestamp    time.Time              `json:"timestamp"`
}

// ButtonPressEvent represents a button press event with timing
type ButtonPressEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency_ms"` // milliseconds from question shown
}

// NewGameEvent creates a new game event
func NewGameEvent(gameID uuid.UUID, eventType EventType) *GameEvent {
	return &GameEvent{
		ID:        uuid.New(),
		GameID:    gameID,
		EventType: eventType,
		Timestamp: time.Now(),
		Data:      make(map[string]interface{}),
	}
}

// WithUser adds user information to the event
func (e *GameEvent) WithUser(userID uuid.UUID) *GameEvent {
	e.UserID = &userID
	return e
}

// WithRound adds round information to the event
func (e *GameEvent) WithRound(roundNumber int) *GameEvent {
	e.RoundNumber = &roundNumber
	return e
}

// WithQuestion adds question information to the event
func (e *GameEvent) WithQuestion(questionID string) *GameEvent {
	e.QuestionID = &questionID
	return e
}

// WithData adds additional data to the event
func (e *GameEvent) WithData(key string, value interface{}) *GameEvent {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data[key] = value
	return e
}

