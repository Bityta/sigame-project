package game

import (
	"sort"

	"github.com/sigame/game/internal/domain"
)

// CalculateScores calculates player scores and ranks
func CalculateScores(players map[string]*domain.Player) []domain.PlayerScore {
	scores := make([]domain.PlayerScore, 0, len(players))

	for _, player := range players {
		scores = append(scores, domain.PlayerScore{
			UserID:   player.UserID,
			Username: player.Username,
			Score:    player.Score,
		})
	}

	// Sort by score descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Assign ranks
	for i := range scores {
		scores[i].Rank = i + 1
	}

	return scores
}

// GetWinners returns the players with the highest score
func GetWinners(scores []domain.PlayerScore) []domain.PlayerScore {
	if len(scores) == 0 {
		return []domain.PlayerScore{}
	}

	maxScore := scores[0].Score
	winners := make([]domain.PlayerScore, 0)

	for _, score := range scores {
		if score.Score == maxScore {
			winners = append(winners, score)
		} else {
			break
		}
	}

	return winners
}

