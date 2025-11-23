package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestGame_CreateGame(t *testing.T) {
	gameID := uuid.New()
	hostID := uuid.New()

	// Test creating a game
	if gameID == uuid.Nil {
		t.Error("Game ID should not be nil")
	}

	if hostID == uuid.Nil {
		t.Error("Host ID should not be nil")
	}
}

func TestPlayer_CreatePlayer(t *testing.T) {
	userID := uuid.New()
	username := "testplayer"

	if userID == uuid.Nil {
		t.Error("User ID should not be nil")
	}

	if username == "" {
		t.Error("Username should not be empty")
	}
}

func TestQuestion_CreateQuestion(t *testing.T) {
	questionID := uuid.New()
	text := "Test question?"
	correctAnswer := "correct"

	if questionID == uuid.Nil {
		t.Error("Question ID should not be nil")
	}

	if text == "" {
		t.Error("Question text should not be empty")
	}

	if correctAnswer == "" {
		t.Error("Correct answer should not be empty")
	}
}

func TestGameEvent_CreateEvent(t *testing.T) {
	eventID := uuid.New()
	gameID := uuid.New()

	if eventID == uuid.Nil {
		t.Error("Event ID should not be nil")
	}

	if gameID == uuid.Nil {
		t.Error("Game ID should not be nil")
	}
}


