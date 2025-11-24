package domain

import (
	"time"

	"github.com/google/uuid"
)

// PlayerRole represents the role of a player in the game
type PlayerRole string

// Player role constants
const (
	PlayerRoleHost   PlayerRole = "host"
	PlayerRolePlayer PlayerRole = "player"
)

// Player represents a player in the game (Entity - internal model)
type Player struct {
	UserID    uuid.UUID  `json:"user_id"`
	Username  string     `json:"username"`
	Role      PlayerRole `json:"role"`
	Score     int        `json:"score"`
	IsActive  bool       `json:"is_active"`
	IsReady   bool       `json:"is_ready"`
	JoinedAt  time.Time  `json:"joined_at"`
	LeftAt    *time.Time `json:"left_at,omitempty"`
}

// PlayerState represents player state for broadcasting (DTO)
// All fields are required
type PlayerState struct {
	UserID   uuid.UUID  `json:"user_id" binding:"required"`
	Username string     `json:"username" binding:"required"`
	Role     PlayerRole `json:"role" binding:"required"`
	Score    int        `json:"score" binding:"required"`
	IsActive bool       `json:"is_active" binding:"required"`
	IsReady  bool       `json:"is_ready" binding:"required"`
}

// PlayerScore represents a player's score entry (DTO)
// All fields are required
type PlayerScore struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	Username string    `json:"username" binding:"required"`
	Score    int       `json:"score" binding:"required"`
	Rank     int       `json:"rank" binding:"required"`
}

// NewPlayer creates a new player
func NewPlayer(userID uuid.UUID, username string, role PlayerRole) *Player {
	return &Player{
		UserID:   userID,
		Username: username,
		Role:     role,
		Score:    0,
		IsActive: true,
		IsReady:  false,
		JoinedAt: time.Now(),
	}
}

// ToState converts Player to PlayerState for broadcasting
func (p *Player) ToState() PlayerState {
	return PlayerState{
		UserID:   p.UserID,
		Username: p.Username,
		Role:     p.Role,
		Score:    p.Score,
		IsActive: p.IsActive,
		IsReady:  p.IsReady,
	}
}

// AddScore adds points to player's score
func (p *Player) AddScore(points int) {
	p.Score += points
}

// SubtractScore subtracts points from player's score
func (p *Player) SubtractScore(points int) {
	p.Score -= points
	if p.Score < 0 {
		p.Score = 0
	}
}

// SetReady marks player as ready
func (p *Player) SetReady(ready bool) {
	p.IsReady = ready
}

// Leave marks player as left
func (p *Player) Leave() {
	p.IsActive = false
	now := time.Now()
	p.LeftAt = &now
}

