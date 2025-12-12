package pack

import "github.com/google/uuid"

type Pack struct {
	ID          uuid.UUID
	Name        string
	Author      string
	Description string
	Rounds      []*Round
}

func (p *Pack) TotalQuestions() int {
	total := 0
	for _, round := range p.Rounds {
		total += round.TotalQuestions()
	}
	return total
}

func (p *Pack) GetRound(roundNumber int) *Round {
	if roundNumber <= 0 || roundNumber > len(p.Rounds) {
		return nil
	}
	return p.Rounds[roundNumber-1]
}

func (p *Pack) TotalRounds() int {
	return len(p.Rounds)
}

