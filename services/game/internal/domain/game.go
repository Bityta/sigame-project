package domain

import (
	"time"

	"github.com/google/uuid"
)

// GameStatus represents the current status of a game
type GameStatus string

// Game status constants
const (
	GameStatusWaiting        GameStatus = "waiting"
	GameStatusRoundStart     GameStatus = "round_start"
	GameStatusQuestionSelect GameStatus = "question_select"
	GameStatusQuestionShow   GameStatus = "question_show"
	GameStatusButtonPress    GameStatus = "button_press"
	GameStatusAnswering      GameStatus = "answering"
	GameStatusAnswerJudging  GameStatus = "answer_judging"
	GameStatusRoundEnd       GameStatus = "round_end"
	GameStatusGameEnd        GameStatus = "game_end"
	GameStatusFinished       GameStatus = "finished"
	GameStatusCancelled      GameStatus = "cancelled"
)

// Game represents the entire game session (Entity - internal model)
type Game struct {
	ID            uuid.UUID              `json:"id"`
	RoomID        uuid.UUID              `json:"room_id"`
	PackID        uuid.UUID              `json:"pack_id"`
	Status        GameStatus             `json:"status"`
	Players       map[uuid.UUID]*Player  `json:"players"`
	Rounds        []*Round               `json:"rounds"`
	CurrentRound  int                    `json:"current_round"`
	CurrentPhase  GameStatus             `json:"current_phase"`
	ActivePlayer  *uuid.UUID             `json:"active_player,omitempty"`
	CurrentTheme  *string                `json:"current_theme,omitempty"`
	CurrentQuestion *Question            `json:"current_question,omitempty"`
	Settings      GameSettings           `json:"settings"`
	StartedAt     *time.Time             `json:"started_at,omitempty"`
	FinishedAt    *time.Time             `json:"finished_at,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// GameSettings holds game configuration (DTO)
// All fields are required
type GameSettings struct {
	TimeForAnswer    int  `json:"time_for_answer" binding:"required"`    // seconds
	TimeForChoice    int  `json:"time_for_choice" binding:"required"`    // seconds
	AllowWrongAnswer bool `json:"allow_wrong_answer" binding:"required"`
	ShowRightAnswer  bool `json:"show_right_answer" binding:"required"`
}

// GameState represents the current state for broadcasting (DTO)
// Required: game_id, status, current_round, players
// Optional: round_name, themes, active_player, current_question, time_remaining, message
type GameState struct {
	GameID          uuid.UUID              `json:"game_id" binding:"required"`
	Status          GameStatus             `json:"status" binding:"required"`
	CurrentRound    int                    `json:"current_round" binding:"required"`
	RoundName       string                 `json:"round_name,omitempty"`
	Themes          []ThemeState           `json:"themes,omitempty"`
	Players         []PlayerState          `json:"players" binding:"required"`
	ActivePlayer    *uuid.UUID             `json:"active_player,omitempty"`
	CurrentQuestion *QuestionState         `json:"current_question,omitempty"`
	TimeRemaining   int                    `json:"time_remaining,omitempty"` // seconds
	Message         string                 `json:"message,omitempty"`
}

// ThemeState represents a theme with questions availability (DTO)
// All fields are required
type ThemeState struct {
	Name      string          `json:"name" binding:"required"`
	Questions []QuestionState `json:"questions" binding:"required"`
}

// QuestionState represents question availability in UI (DTO)
// Required: id, price, available
// Optional: text, media_type (only when question is shown)
type QuestionState struct {
	ID        string `json:"id" binding:"required"`
	Price     int    `json:"price" binding:"required"`
	Available bool   `json:"available" binding:"required"`
	Text      string `json:"text,omitempty"`      // only when shown
	MediaType string `json:"media_type,omitempty"` // only when shown
}

// CreateGameRequest is the request to create a new game (DTO)
// All fields are required
type CreateGameRequest struct {
	RoomID   uuid.UUID    `json:"room_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	PackID   uuid.UUID    `json:"pack_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	Players  []PlayerInfo `json:"players" binding:"required,min=2"`
	Settings GameSettings `json:"settings" binding:"required"`
}

// PlayerInfo is player information for game creation (DTO)
// All fields are required
type PlayerInfo struct {
	UserID   uuid.UUID `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440002"`
	Username string    `json:"username" binding:"required" example:"player1"`
	Role     string    `json:"role" binding:"required" example:"player"`
}

// CreateGameResponse is the response after creating a game (DTO)
// All fields are required
type CreateGameResponse struct {
	GameID       uuid.UUID `json:"game_id" binding:"required"`
	WebSocketURL string    `json:"websocket_url" binding:"required"`
	Status       string    `json:"status" binding:"required"`
}

// GetGameResponse is the response for getting game info (DTO)
// Required: game_id, room_id, pack_id, status, current_round, players
// Optional: started_at, finished_at
type GetGameResponse struct {
	GameID       uuid.UUID     `json:"game_id" binding:"required"`
	RoomID       uuid.UUID     `json:"room_id" binding:"required"`
	PackID       uuid.UUID     `json:"pack_id" binding:"required"`
	Status       GameStatus    `json:"status" binding:"required"`
	CurrentRound int           `json:"current_round" binding:"required"`
	Players      []PlayerState `json:"players" binding:"required"`
	StartedAt    *time.Time    `json:"started_at,omitempty"`
	FinishedAt   *time.Time    `json:"finished_at,omitempty"`
}

