package pack

import (
	"fmt"

	"github.com/google/uuid"
	domainPack "sigame/game/internal/domain/pack"
)

func convertPackResponse(packResp *PackContentResponse, packID uuid.UUID) (*domainPack.Pack, error) {
	if packResp.Error != "" {
		return nil, fmt.Errorf("pack service returned error: %s", packResp.Error)
	}

	p := &domainPack.Pack{
		ID:     packID,
		Name:   packResp.Name,
		Author: packResp.Author,
		Rounds: make([]*domainPack.Round, len(packResp.Rounds)),
	}

	for i, r := range packResp.Rounds {
		p.Rounds[i] = convertRound(r)
	}

	return p, nil
}

func convertRound(r RoundJSON) *domainPack.Round {
	round := &domainPack.Round{
		ID:          r.ID,
		RoundNumber: r.RoundNumber,
		Name:        r.Name,
		Themes:      make([]*domainPack.Theme, len(r.Themes)),
	}

	for j, t := range r.Themes {
		round.Themes[j] = convertTheme(t)
	}

	return round
}

func convertTheme(t ThemeJSON) *domainPack.Theme{
	theme := &domainPack.Theme{
		ID:        t.ID,
		Name:      t.Name,
		Questions: make([]*domainPack.Question, len(t.Questions)),
	}

	for k, q := range t.Questions {
		theme.Questions[k] = convertQuestion(q)
	}

	return theme
}

func convertQuestion(q QuestionJSON) *domainPack.Question {
	return &domainPack.Question{
		ID:              q.ID,
		Price:           q.Price,
		Text:            q.Text,
		Answer:          q.Answer,
		MediaType:       q.MediaType,
		MediaURL:        q.MediaURL,
		MediaDurationMs: q.MediaDurationMs,
		Used:            false,
	}
}
