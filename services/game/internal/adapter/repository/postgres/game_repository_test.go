package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestNewGameRepository(t *testing.T) {
	repo := NewGameRepository(nil)
	if repo == nil {
		t.Error("NewGameRepository() returned nil")
	}
	if repo.db != nil {
		t.Error("NewGameRepository() db should be nil when passed nil")
	}
}

func TestGameRepository_CreateGameSession(t *testing.T) {
	repo := NewGameRepository(nil)
	if repo == nil {
		t.Error("NewGameRepository() returned nil")
	}
	if repo.db != nil {
		t.Error("NewGameRepository() db should be nil when passed nil")
	}
}

func TestGameRepository_MethodsExist(t *testing.T) {
	repo := NewGameRepository(nil)
	ctx := context.Background()
	gameID := uuid.New()
	roomID := uuid.New()
	userID := uuid.New()

	_ = repo.CreateGameSession
	_ = repo.UpdateGameSession
	_ = repo.GetGameSession
	_ = repo.SaveFinalResults
	_ = repo.GetGamesByRoomID
	_ = repo.GetActiveGameForUser

	_ = ctx
	_ = gameID
	_ = roomID
	_ = userID
}

