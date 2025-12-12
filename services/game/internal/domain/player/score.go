package player

import "github.com/google/uuid"

type Score struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	Username string    `json:"username" binding:"required"`
	Score    int       `json:"score" binding:"required"`
	Rank     int       `json:"rank" binding:"required"`
}

func NewScore(userID uuid.UUID, username string, score int) Score {
	return Score{
		UserID:   userID,
		Username: username,
		Score:    score,
		Rank:     0,
	}
}

