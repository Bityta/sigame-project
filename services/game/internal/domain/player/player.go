package player

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	UserID      uuid.UUID
	Username    string
	AvatarURL   string
	Role        Role
	Score       int
	IsActive    bool
	IsReady     bool
	IsConnected bool
	JoinedAt    time.Time
	LeftAt      *time.Time
}

type State struct {
	UserID      uuid.UUID `json:"userId" binding:"required"`
	Username    string    `json:"username" binding:"required"`
	AvatarURL   string    `json:"avatarUrl,omitempty"`
	Role        Role      `json:"role" binding:"required"`
	Score       int       `json:"score" binding:"required"`
	IsActive    bool      `json:"isActive" binding:"required"`
	IsReady     bool      `json:"isReady" binding:"required"`
	IsConnected bool      `json:"isConnected" binding:"required"`
}

func New(userID uuid.UUID, username string, avatarURL string, role Role) *Player {
	return &Player{
		UserID:      userID,
		Username:    username,
		AvatarURL:   avatarURL,
		Role:        role,
		Score:       0,
		IsActive:    true,
		IsReady:     false,
		IsConnected: false,
		JoinedAt:    time.Now(),
	}
}

func (p *Player) ToState() State {
	return State{
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

func (p *Player) SetConnected(connected bool) {
	p.IsConnected = connected
}

func (p *Player) AddScore(points int) {
	p.Score += points
}

func (p *Player) SubtractScore(points int) {
	p.Score -= points
	if p.Score < 0 {
		p.Score = 0
	}
}

func (p *Player) SetReady(ready bool) {
	p.IsReady = ready
}

func (p *Player) Leave() {
	p.IsActive = false
	now := time.Now()
	p.LeftAt = &now
}

func (p *Player) IsHost() bool {
	return p.Role == RoleHost
}

func (p *Player) CanPressButton() bool {
	return p.IsActive && !p.IsHost()
}

func (p *Player) CanAnswer() bool {
	return p.IsActive && !p.IsHost()
}

