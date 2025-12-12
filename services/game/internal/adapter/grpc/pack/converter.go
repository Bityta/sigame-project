package grpc

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
	"github.com/sigame/game/internal/domain/event"
)

func convertPackResponse(packResp *PackContentResponse, packID uuid.UUID) (*domain.Pack, error) {
	if packResp.Error != "" {
		return nil, fmt.Errorf("pack service returned error: %s", packResp.Error)
	}

	pack := &domain.Pack{
		ID:     packID,
		Name:   packResp.Name,
		Author: packResp.Author,
		Rounds: make([]*domain.Round, len(packResp.Rounds)),
	}

	for i, r := range packResp.Rounds {
		pack.Rounds[i] = convertRound(r)
	}

	return pack, nil
}

func convertRound(r RoundJSON) *domain.Round {
	round := &domain.Round{
		ID:          r.ID,
		RoundNumber: r.RoundNumber,
		Name:        r.Name,
		Themes:      make([]*domain.Theme, len(r.Themes)),
	}

	for j, t := range r.Themes {
		round.Themes[j] = convertTheme(t)
	}

	return round
}

func convertTheme(t ThemeJSON) *domain.Theme {
	theme := &domain.Theme{
		ID:        t.ID,
		Name:      t.Name,
		Questions: make([]*domain.Question, len(t.Questions)),
	}

	for k, q := range t.Questions {
		theme.Questions[k] = convertQuestion(q)
	}

	return theme
}

func convertQuestion(q QuestionJSON) *domain.Question {
	return &domain.Question{
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

