package event

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	GameID      uuid.UUID
	EventType   Type
	UserID      *uuid.UUID
	RoundNumber *int
	QuestionID  *string
	Data        map[string]interface{}
	Timestamp   time.Time
}

func New(gameID uuid.UUID, eventType Type) *Event {
	return &Event{
		ID:        uuid.New(),
		GameID:    gameID,
		EventType: eventType,
		Timestamp: time.Now(),
		Data:      make(map[string]interface{}),
	}
}

func (e *Event) WithUser(userID uuid.UUID) *Event {
	e.UserID = &userID
	return e
}

func (e *Event) WithRound(roundNumber int) *Event {
	e.RoundNumber = &roundNumber
	return e
}

func (e *Event) WithQuestion(questionID string) *Event {
	e.QuestionID = &questionID
	return e
}

func (e *Event) WithData(key string, value interface{}) *Event {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data[key] = value
	return e
}

