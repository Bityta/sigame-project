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
	UserID      uuid.UUID  `json:"user_id"`
	Username    string     `json:"username"`
	AvatarURL   string     `json:"avatar_url"`
	Role        PlayerRole `json:"role"`
	Score       int        `json:"score"`
	IsActive    bool       `json:"is_active"`
	IsReady     bool       `json:"is_ready"`
	IsConnected bool       `json:"is_connected"`
	JoinedAt    time.Time  `json:"joined_at"`
	LeftAt      *time.Time `json:"left_at,omitempty"`
}

// PlayerState represents player state for broadcasting (DTO)
// All fields are required
type PlayerState struct {
	UserID      uuid.UUID  `json:"userId" binding:"required"`
	Username    string     `json:"username" binding:"required"`
	AvatarURL   string     `json:"avatarUrl,omitempty"`
	Role        PlayerRole `json:"role" binding:"required"`
	Score       int        `json:"score" binding:"required"`
	IsActive    bool       `json:"isActive" binding:"required"`
	IsReady     bool       `json:"isReady" binding:"required"`
	IsConnected bool       `json:"isConnected" binding:"required"`
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
func NewPlayer(userID uuid.UUID, username string, avatarURL string, role PlayerRole) *Player {
	return &Player{
		UserID:      userID,
		Username:    username,
		AvatarURL:   avatarURL,
		Role:        role,
		Score:       0,
		IsActive:    true,
		IsReady:     false,
		IsConnected: false, // Will be set to true when they connect via WebSocket
		JoinedAt:    time.Now(),
	}
}

// ToState converts Player to PlayerState for broadcasting
func (p *Player) ToState() PlayerState {
	return PlayerState{
		UserID:      p.UserID,
		Username:    p.Username,
		AvatarURL:   p.AvatarURL,
		Role:        p.Role,
		Score:       p.Score,
		IsActive:    p.IsActive,
		IsReady:     p.IsReady,
		IsConnected: p.IsConnected,
	}
}

// SetConnected updates the connection status
func (p *Player) SetConnected(connected bool) {
	p.IsConnected = connected
}

// AddScore adds points to player's score
func (p *Player) AddScore(points int) {
	p.Score += points
}

// SubtractScore subtracts points from player's score (can go negative like in real SIGame)
func (p *Player) SubtractScore(points int) {
	p.Score -= points
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

