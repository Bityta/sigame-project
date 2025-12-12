package game

import (
	"sort"

	"sigame/game/internal/core/scoring"
	"sigame/game/internal/domain/player"
)

func (m *Manager) calculateWinners() []player.Score {
	scores := m.calculateFinalScores()

	winners := make([]player.Score, 0)
	for i, score := range scores {
		if i >= TopWinnersCount {
			break
		}
		winners = append(winners, score)
	}

	return winners
}

func (m *Manager) calculateFinalScores() []player.Score {
	scores := make([]player.Score, 0)

	for userID, p := range m.game.Players {
		if p.Role == player.RoleHost {
			continue
		}
		scores = append(scores, player.NewScore(userID, p.Username, p.Score))
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	for i := range scores {
		scores[i].Rank = i + scoring.RankStartIndex
	}

	return scores
}

