package handler

import (
	"github.com/google/uuid"
)

type CreateGameRequest struct {
	RoomID   uuid.UUID    `json:"room_id" binding:"required"`
	PackID   uuid.UUID    `json:"pack_id" binding:"required"`
	Players  []PlayerInfo `json:"players" binding:"required,min=2"`
	Settings GameSettings `json:"settings" binding:"required"`
}

type PlayerInfo struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	Username  string    `json:"username" binding:"required"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role" binding:"required"`
}

type GameSettings struct {
	TimeForAnswer int `json:"time_for_answer" binding:"required"`
	TimeForChoice int `json:"time_for_choice" binding:"required"`
}

type CreateGameResponse struct {
	GameID       uuid.UUID `json:"game_id"`
	WebSocketURL string    `json:"websocket_url"`
	Status       string    `json:"status"`
}

type GetGameResponse struct {
	GameID       uuid.UUID `json:"game_id"`
	RoomID       uuid.UUID `json:"room_id"`
	PackID       uuid.UUID `json:"pack_id"`
	Status       string    `json:"status"`
	CurrentRound int       `json:"current_round"`
}

