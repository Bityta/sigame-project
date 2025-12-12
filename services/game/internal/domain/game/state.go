package game

import (
	"github.com/google/uuid"
	"sigame/game/internal/domain/pack"
	"sigame/game/internal/domain/player"
)

type State struct {
	GameID          uuid.UUID           `json:"gameId" binding:"required"`
	Status          Status              `json:"status" binding:"required"`
	CurrentRound    int                 `json:"currentRound" binding:"required"`
	RoundName       string              `json:"roundName,omitempty"`
	Themes          []pack.ThemeState   `json:"themes,omitempty"`
	Players         []player.State      `json:"players" binding:"required"`
	ActivePlayer    *uuid.UUID          `json:"activePlayer,omitempty"`
	CurrentQuestion *pack.QuestionState `json:"currentQuestion,omitempty"`
	TimeRemaining   int                 `json:"timeRemaining,omitempty"`
	Message         string              `json:"message,omitempty"`
	AllRounds       []RoundOverview     `json:"allRounds,omitempty"`
	Winners         []player.Score      `json:"winners,omitempty"`
	FinalScores     []player.Score      `json:"finalScores,omitempty"`

	StakeInfo     *StakeInfo         `json:"stakeInfo,omitempty"`
	SecretTarget  *uuid.UUID         `json:"secretTarget,omitempty"`
	ForAllResults []ForAllResult     `json:"forAllResults,omitempty"`
}

type RoundOverview struct {
	RoundNumber int      `json:"roundNumber" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	ThemeNames  []string `json:"themeNames" binding:"required"`
}

type StakeInfo struct {
	MinBet     int  `json:"minBet"`
	MaxBet     int  `json:"maxBet"`
	CurrentBet int  `json:"currentBet"`
	IsAllIn    bool `json:"isAllIn"`
}

type ForAllResult struct {
	UserID     uuid.UUID `json:"userId"`
	Username   string    `json:"username"`
	Answer     string    `json:"answer"`
	IsCorrect  bool      `json:"isCorrect"`
	ScoreDelta int       `json:"scoreDelta"`
}

