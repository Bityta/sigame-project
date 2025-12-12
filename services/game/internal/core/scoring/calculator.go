package scoring

import (
	"sort"

	"sigame/game/internal/domain/game"
	"sigame/game/internal/domain/pack"
	"sigame/game/internal/domain/player"
)

func CalculateScores(players map[string]*domain.Player) []domain.PlayerScore {
	scores := make([]domain.PlayerScore, 0, len(players))

	for _, player := range players {
		scores = append(scores, domain.PlayerScore{
			UserID:   player.UserID,
			Username: player.Username,
			Score:    player.Score,
		})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	for i := range scores {
		scores[i].Rank = i + RankStartIndex
	}

	return scores
}

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

