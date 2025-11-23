package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGame_AddPlayer(t *testing.T) {
	game := &Game{
		ID:      uuid.New(),
		Players: make(map[uuid.UUID]*Player),
	}

	player := &Player{
		ID:       uuid.New(),
		Username: "testplayer",
	}

	game.AddPlayer(player)

	if len(game.Players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(game.Players))
	}

	if game.Players[player.ID] != player {
		t.Error("Player not added correctly")
	}
}

func TestGame_RemovePlayer(t *testing.T) {
	playerID := uuid.New()
	game := &Game{
		ID: uuid.New(),
		Players: map[uuid.UUID]*Player{
			playerID: {
				ID:       playerID,
				Username: "testplayer",
			},
		},
	}

	game.RemovePlayer(playerID)

	if len(game.Players) != 0 {
		t.Errorf("Expected 0 players, got %d", len(game.Players))
	}
}

func TestPlayer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		player  Player
		wantErr bool
	}{
		{
			name: "valid player",
			player: Player{
				ID:       uuid.New(),
				Username: "testplayer",
			},
			wantErr: false,
		},
		{
			name: "empty username",
			player: Player{
				ID:       uuid.New(),
				Username: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.player.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Player.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestion_IsCorrect(t *testing.T) {
	question := &Question{
		ID:            uuid.New(),
		Text:          "Test question?",
		CorrectAnswer: "correct",
	}

	tests := []struct {
		name   string
		answer string
		want   bool
	}{
		{
			name:   "correct answer",
			answer: "correct",
			want:   true,
		},
		{
			name:   "incorrect answer",
			answer: "wrong",
			want:   false,
		},
		{
			name:   "case insensitive correct",
			answer: "CORRECT",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := question.IsCorrect(tt.answer); got != tt.want {
				t.Errorf("Question.IsCorrect() = %v, want %v", got, tt.want)
			}
		})
	}
}

