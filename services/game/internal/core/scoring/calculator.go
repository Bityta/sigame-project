package scoring

import (
	"sort"

	"sigame/game/internal/domain/player"
)

func CalculateScores(players map[string]*player.Player) []player.Score {
	scores := make([]player.Score, 0, len(players))

	for _, p := range players {
		scores = append(scores, player.Score{
			UserID:   p.UserID,
			Username: p.Username,
			Score:    p.Score,
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

func GetWinners(scores []player.Score) []player.Score {
	if len(scores) == 0 {
		return []player.Score{}
	}

	maxScore := scores[0].Score
	winners := make([]player.Score, 0)

	for _, score := range scores {
		if score.Score == maxScore {
			winners = append(winners, score)
		} else {
			break
		}
	}

	return winners
}

