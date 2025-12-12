package scoring

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
)

func TestCalculateScores(t *testing.T) {
	player1 := &domain.Player{
		UserID:   uuid.New(),
		Username: "player1",
		Score:    100,
	}

	player2 := &domain.Player{
		UserID:   uuid.New(),
		Username: "player2",
		Score:    200,
	}

	player3 := &domain.Player{
		UserID:   uuid.New(),
		Username: "player3",
		Score:    150,
	}

	players := map[string]*domain.Player{
		"player1": player1,
		"player2": player2,
		"player3": player3,
	}

	scores := CalculateScores(players)

	if len(scores) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(scores))
	}

	if scores[0].Score != 200 {
		t.Errorf("Expected highest score 200, got %d", scores[0].Score)
	}

	if scores[0].Rank != RankStartIndex {
		t.Errorf("Expected rank %d for highest score, got %d", RankStartIndex, scores[0].Rank)
	}

	if scores[1].Score != 150 {
		t.Errorf("Expected second score 150, got %d", scores[1].Score)
	}

	if scores[1].Rank != RankStartIndex+1 {
		t.Errorf("Expected rank %d for second score, got %d", RankStartIndex+1, scores[1].Rank)
	}

	if scores[2].Score != 100 {
		t.Errorf("Expected third score 100, got %d", scores[2].Score)
	}

	if scores[2].Rank != RankStartIndex+2 {
		t.Errorf("Expected rank %d for third score, got %d", RankStartIndex+2, scores[2].Rank)
	}

	if scores[0].UserID != player2.UserID {
		t.Error("Expected player2 to have highest score")
	}

	if scores[1].UserID != player3.UserID {
		t.Error("Expected player3 to have second score")
	}

	if scores[2].UserID != player1.UserID {
		t.Error("Expected player1 to have third score")
	}
}

func TestCalculateScores_Empty(t *testing.T) {
	players := map[string]*domain.Player{}

	scores := CalculateScores(players)

	if len(scores) != 0 {
		t.Errorf("Expected 0 scores for empty players map, got %d", len(scores))
	}
}

func TestCalculateScores_Sorted(t *testing.T) {
	players := map[string]*domain.Player{
		"p1": {UserID: uuid.New(), Username: "p1", Score: 50},
		"p2": {UserID: uuid.New(), Username: "p2", Score: 100},
		"p3": {UserID: uuid.New(), Username: "p3", Score: 75},
		"p4": {UserID: uuid.New(), Username: "p4", Score: 200},
		"p5": {UserID: uuid.New(), Username: "p5", Score: 25},
	}

	scores := CalculateScores(players)

	if len(scores) != 5 {
		t.Fatalf("Expected 5 scores, got %d", len(scores))
	}

	for i := 0; i < len(scores)-1; i++ {
		if scores[i].Score < scores[i+1].Score {
			t.Errorf("Scores not sorted descending: scores[%d]=%d < scores[%d]=%d",
				i, scores[i].Score, i+1, scores[i+1].Score)
		}
	}

	expectedScores := []int{200, 100, 75, 50, 25}
	for i, expected := range expectedScores {
		if scores[i].Score != expected {
			t.Errorf("Expected score[%d]=%d, got %d", i, expected, scores[i].Score)
		}
	}

	for i, score := range scores {
		expectedRank := i + RankStartIndex
		if score.Rank != expectedRank {
			t.Errorf("Expected rank %d for score[%d], got %d", expectedRank, i, score.Rank)
		}
	}
}

func TestGetWinners(t *testing.T) {
	scores := []domain.PlayerScore{
		{UserID: uuid.New(), Username: "winner1", Score: 200, Rank: 1},
		{UserID: uuid.New(), Username: "winner2", Score: 200, Rank: 1},
		{UserID: uuid.New(), Username: "loser1", Score: 100, Rank: 3},
		{UserID: uuid.New(), Username: "loser2", Score: 50, Rank: 4},
	}

	winners := GetWinners(scores)

	if len(winners) != 2 {
		t.Errorf("Expected 2 winners, got %d", len(winners))
	}

	for _, winner := range winners {
		if winner.Score != 200 {
			t.Errorf("Expected winner score 200, got %d", winner.Score)
		}
	}

	if winners[0].Username != "winner1" {
		t.Errorf("Expected first winner username 'winner1', got '%s'", winners[0].Username)
	}

	if winners[1].Username != "winner2" {
		t.Errorf("Expected second winner username 'winner2', got '%s'", winners[1].Username)
	}
}

func TestGetWinners_Empty(t *testing.T) {
	scores := []domain.PlayerScore{}

	winners := GetWinners(scores)

	if len(winners) != 0 {
		t.Errorf("Expected 0 winners for empty scores, got %d", len(winners))
	}

	if winners == nil {
		t.Error("GetWinners should return empty slice, not nil")
	}
}

func TestGetWinners_MultipleWinners(t *testing.T) {
	scores := []domain.PlayerScore{
		{UserID: uuid.New(), Username: "winner1", Score: 150, Rank: 1},
		{UserID: uuid.New(), Username: "winner2", Score: 150, Rank: 1},
		{UserID: uuid.New(), Username: "winner3", Score: 150, Rank: 1},
		{UserID: uuid.New(), Username: "loser", Score: 100, Rank: 4},
	}

	winners := GetWinners(scores)

	if len(winners) != 3 {
		t.Errorf("Expected 3 winners, got %d", len(winners))
	}

	for i, winner := range winners {
		if winner.Score != 150 {
			t.Errorf("Expected winner[%d] score 150, got %d", i, winner.Score)
		}
	}
}

func TestGetWinners_SingleWinner(t *testing.T) {
	scores := []domain.PlayerScore{
		{UserID: uuid.New(), Username: "winner", Score: 200, Rank: 1},
		{UserID: uuid.New(), Username: "loser1", Score: 100, Rank: 2},
		{UserID: uuid.New(), Username: "loser2", Score: 50, Rank: 3},
	}

	winners := GetWinners(scores)

	if len(winners) != 1 {
		t.Errorf("Expected 1 winner, got %d", len(winners))
	}

	if winners[0].Score != 200 {
		t.Errorf("Expected winner score 200, got %d", winners[0].Score)
	}

	if winners[0].Username != "winner" {
		t.Errorf("Expected winner username 'winner', got '%s'", winners[0].Username)
	}
}

func TestGetWinners_AllSameScore(t *testing.T) {
	scores := []domain.PlayerScore{
		{UserID: uuid.New(), Username: "p1", Score: 100, Rank: 1},
		{UserID: uuid.New(), Username: "p2", Score: 100, Rank: 1},
		{UserID: uuid.New(), Username: "p3", Score: 100, Rank: 1},
	}

	winners := GetWinners(scores)

	if len(winners) != 3 {
		t.Errorf("Expected 3 winners when all have same score, got %d", len(winners))
	}
}

func TestCalculateScores_WithZeroScores(t *testing.T) {
	players := map[string]*domain.Player{
		"p1": {UserID: uuid.New(), Username: "p1", Score: 0},
		"p2": {UserID: uuid.New(), Username: "p2", Score: 100},
		"p3": {UserID: uuid.New(), Username: "p3", Score: 0},
	}

	scores := CalculateScores(players)

	if len(scores) != 3 {
		t.Fatalf("Expected 3 scores, got %d", len(scores))
	}

	if scores[0].Score != 100 {
		t.Errorf("Expected highest score 100, got %d", scores[0].Score)
	}

	if scores[1].Score != 0 {
		t.Errorf("Expected second score 0, got %d", scores[1].Score)
	}

	if scores[2].Score != 0 {
		t.Errorf("Expected third score 0, got %d", scores[2].Score)
	}
}

