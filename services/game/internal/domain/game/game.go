package game

import (
	"time"

	"github.com/google/uuid"
	"sigame/game/internal/domain/pack"
	"sigame/game/internal/domain/player"
)

type Game struct {
	ID              uuid.UUID
	RoomID          uuid.UUID
	PackID          uuid.UUID
	Status          Status
	Players         map[uuid.UUID]*player.Player
	Rounds          []*pack.Round
	CurrentRound    int
	CurrentPhase    Status
	ActivePlayer    *uuid.UUID
	CurrentTheme    *string
	CurrentQuestion *pack.Question
	Settings        Settings
	Winners         []player.Score
	FinalScores     []player.Score
	StartedAt       *time.Time
	FinishedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func New(roomID, packID uuid.UUID, settings Settings, rounds []*pack.Round) *Game {
	now := time.Now()
	return &Game{
		ID:           uuid.New(),
		RoomID:       roomID,
		PackID:       packID,
		Status:       StatusWaiting,
		Players:      make(map[uuid.UUID]*player.Player),
		Rounds:       rounds,
		CurrentRound: 0,
		CurrentPhase: StatusWaiting,
		Settings:     settings,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (g *Game) AddPlayer(p *player.Player) error {
	if _, exists := g.Players[p.UserID]; exists {
		return ErrPlayerAlreadyExists
	}
	g.Players[p.UserID] = p
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Game) GetPlayer(userID uuid.UUID) (*player.Player, error) {
	p, exists := g.Players[userID]
	if !exists {
		return nil, ErrPlayerNotFound
	}
	return p, nil
}

func (g *Game) Start() error {
	if g.Status != StatusWaiting {
		return ErrInvalidGameStatus
	}
	
	now := time.Now()
	g.StartedAt = &now
	g.Status = StatusRoundsOverview
	g.CurrentPhase = StatusRoundsOverview
	g.UpdatedAt = now
	
	return nil
}

func (g *Game) Finish() {
	now := time.Now()
	g.Status = StatusFinished
	g.CurrentPhase = StatusFinished
	g.FinishedAt = &now
	g.UpdatedAt = now
}

func (g *Game) Cancel() {
	g.Status = StatusCancelled
	g.CurrentPhase = StatusCancelled
	g.UpdatedAt = time.Now()
}

func (g *Game) UpdateStatus(status Status) {
	g.Status = status
	g.CurrentPhase = status
	g.UpdatedAt = time.Now()
}

func (g *Game) SetActivePlayer(userID uuid.UUID) {
	g.ActivePlayer = &userID
}

func (g *Game) ClearActivePlayer() {
	g.ActivePlayer = nil
}

func (g *Game) GetCurrentRound() (*pack.Round, error) {
	if g.CurrentRound <= 0 || g.CurrentRound > len(g.Rounds) {
		return nil, ErrInvalidRound
	}
	return g.Rounds[g.CurrentRound-1], nil
}

func (g *Game) SetCurrentQuestion(q *pack.Question, themeName string) {
	g.CurrentQuestion = q
	g.CurrentTheme = &themeName
}

func (g *Game) ClearCurrentQuestion() {
	g.CurrentQuestion = nil
	g.CurrentTheme = nil
}

func (g *Game) IsFinished() bool {
	return g.Status == StatusFinished || g.Status == StatusCancelled
}

func (g *Game) GetActivePlayers() []*player.Player {
	active := make([]*player.Player, 0)
	for _, p := range g.Players {
		if p.IsActive {
			active = append(active, p)
		}
	}
	return active
}

func (g *Game) GetHost() (*player.Player, error) {
	for _, p := range g.Players {
		if p.Role == player.RoleHost {
			return p, nil
		}
	}
	return nil, ErrHostNotFound
}

