package domain

import (
	"time"

	"github.com/google/uuid"
)

// GameStatus represents the current status of a game
type GameStatus string

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

// Game represents the entire game session
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

// GameSettings holds game configuration
type GameSettings struct {
	TimeForAnswer    int  `json:"time_for_answer"`    // seconds
	TimeForChoice    int  `json:"time_for_choice"`    // seconds
	AllowWrongAnswer bool `json:"allow_wrong_answer"`
	ShowRightAnswer  bool `json:"show_right_answer"`
}

// GameState represents the current state for broadcasting
type GameState struct {
	GameID          uuid.UUID              `json:"game_id"`
	Status          GameStatus             `json:"status"`
	CurrentRound    int                    `json:"current_round"`
	RoundName       string                 `json:"round_name,omitempty"`
	Themes          []ThemeState           `json:"themes,omitempty"`
	Players         []PlayerState          `json:"players"`
	ActivePlayer    *uuid.UUID             `json:"active_player,omitempty"`
	CurrentQuestion *QuestionState         `json:"current_question,omitempty"`
	TimeRemaining   int                    `json:"time_remaining,omitempty"` // seconds
	Message         string                 `json:"message,omitempty"`
}

// ThemeState represents a theme with questions availability
type ThemeState struct {
	Name      string          `json:"name"`
	Questions []QuestionState `json:"questions"`
}

// QuestionState represents question availability in UI
type QuestionState struct {
	ID        string `json:"id"`
	Price     int    `json:"price"`
	Available bool   `json:"available"`
	Text      string `json:"text,omitempty"`      // only when shown
	MediaType string `json:"media_type,omitempty"` // only when shown
}

// CreateGameRequest is the request to create a new game
type CreateGameRequest struct {
	RoomID   uuid.UUID    `json:"room_id" binding:"required"`
	PackID   uuid.UUID    `json:"pack_id" binding:"required"`
	Players  []PlayerInfo `json:"players" binding:"required,min=2"`
	Settings GameSettings `json:"settings"`
}

// PlayerInfo is player information for game creation
type PlayerInfo struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	Username string    `json:"username" binding:"required"`
	Role     string    `json:"role" binding:"required"`
}

// CreateGameResponse is the response after creating a game
type CreateGameResponse struct {
	GameID       uuid.UUID `json:"game_id"`
	WebSocketURL string    `json:"websocket_url"`
	Status       string    `json:"status"`
}

// GetGameResponse is the response for getting game info
type GetGameResponse struct {
	GameID       uuid.UUID     `json:"game_id"`
	RoomID       uuid.UUID     `json:"room_id"`
	PackID       uuid.UUID     `json:"pack_id"`
	Status       GameStatus    `json:"status"`
	CurrentRound int           `json:"current_round"`
	Players      []PlayerState `json:"players"`
	StartedAt    *time.Time    `json:"started_at,omitempty"`
	FinishedAt   *time.Time    `json:"finished_at,omitempty"`
}

