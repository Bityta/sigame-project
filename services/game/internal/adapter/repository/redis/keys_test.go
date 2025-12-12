package redis

import (
	"testing"

	"github.com/google/uuid"
)

func TestPackKey(t *testing.T) {
	packID := uuid.New()
	key := packKey(packID)
	expected := "pack:" + packID.String() + ":content"

	if key != expected {
		t.Errorf("packKey() = %v, want %v", key, expected)
	}
}

func TestGameStateKey(t *testing.T) {
	gameID := uuid.New()
	key := gameStateKey(gameID)
	expected := "game:" + gameID.String() + ":state"

	if key != expected {
		t.Errorf("gameStateKey() = %v, want %v", key, expected)
	}
}

func TestGameScoresKey(t *testing.T) {
	gameID := uuid.New()
	key := gameScoresKey(gameID)
	expected := "game:" + gameID.String() + ":scores"

	if key != expected {
		t.Errorf("gameScoresKey() = %v, want %v", key, expected)
	}
}

func TestGamePlayersKey(t *testing.T) {
	gameID := uuid.New()
	key := gamePlayersKey(gameID)
	expected := "game:" + gameID.String() + ":players"

	if key != expected {
		t.Errorf("gamePlayersKey() = %v, want %v", key, expected)
	}
}

func TestGameMetadataKey(t *testing.T) {
	gameID := uuid.New()
	key := gameMetadataKey(gameID)
	expected := "game:" + gameID.String() + ":meta"

	if key != expected {
		t.Errorf("gameMetadataKey() = %v, want %v", key, expected)
	}
}

func TestActiveGamesKey(t *testing.T) {
	key := activeGamesKey()
	expected := "games:active"

	if key != expected {
		t.Errorf("activeGamesKey() = %v, want %v", key, expected)
	}
}

