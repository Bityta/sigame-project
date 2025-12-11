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
	GameStatusRoundsOverview GameStatus = "rounds_overview"
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

	// Special question type statuses
	GameStatusSecretTransfer  GameStatus = "secret_transfer"   // Хост выбирает игрока для передачи кота
	GameStatusStakeBetting    GameStatus = "stake_betting"     // Игрок делает ставку
	GameStatusForAllAnswering GameStatus = "for_all_answering" // Все отвечают одновременно
	GameStatusForAllResults   GameStatus = "for_all_results"   // Показ результатов вопроса для всех
)

// Game represents the entire game session (Entity - internal model)
type Game struct {
	ID              uuid.UUID              `json:"id"`
	RoomID          uuid.UUID              `json:"room_id"`
	PackID          uuid.UUID              `json:"pack_id"`
	Status          GameStatus             `json:"status"`
	Players         map[uuid.UUID]*Player  `json:"players"`
	Rounds          []*Round               `json:"rounds"`
	CurrentRound    int                    `json:"current_round"`
	CurrentPhase    GameStatus             `json:"current_phase"`
	ActivePlayer    *uuid.UUID             `json:"active_player,omitempty"`
	CurrentTheme    *string                `json:"current_theme,omitempty"`
	CurrentQuestion *Question              `json:"current_question,omitempty"`
	Settings        GameSettings           `json:"settings"`
	Winners         []PlayerScore          `json:"winners,omitempty"`
	FinalScores     []PlayerScore          `json:"final_scores,omitempty"`
	StartedAt       *time.Time             `json:"started_at,omitempty"`
	FinishedAt      *time.Time             `json:"finished_at,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// GameSettings holds game configuration (DTO)
// All fields are required
type GameSettings struct {
	TimeForAnswer int `json:"time_for_answer" binding:"required"` // seconds
	TimeForChoice int `json:"time_for_choice" binding:"required"` // seconds
}

// RoundOverview represents a round summary for the overview screen
type RoundOverview struct {
	RoundNumber int      `json:"roundNumber" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	ThemeNames  []string `json:"themeNames" binding:"required"`
}

// GameState represents the current state for broadcasting (DTO)
// Required: gameId, status, currentRound, players
// Optional: roundName, themes, activePlayer, currentQuestion, timeRemaining, message, allRounds, winners, finalScores
type GameState struct {
	GameID          uuid.UUID              `json:"gameId" binding:"required"`
	Status          GameStatus             `json:"status" binding:"required"`
	CurrentRound    int                    `json:"currentRound" binding:"required"`
	RoundName       string                 `json:"roundName,omitempty"`
	Themes          []ThemeState           `json:"themes,omitempty"`
	Players         []PlayerState          `json:"players" binding:"required"`
	ActivePlayer    *uuid.UUID             `json:"activePlayer,omitempty"`
	CurrentQuestion *QuestionState         `json:"currentQuestion,omitempty"`
	TimeRemaining   int                    `json:"timeRemaining,omitempty"` // seconds
	Message         string                 `json:"message,omitempty"`
	AllRounds       []RoundOverview        `json:"allRounds,omitempty"` // for rounds_overview status
	Winners         []PlayerScore          `json:"winners,omitempty"`   // for game_end status
	FinalScores     []PlayerScore          `json:"finalScores,omitempty"` // for game_end status

	// Special question type fields
	StakeInfo       *StakeInfo             `json:"stakeInfo,omitempty"`       // for stake_betting status
	SecretTarget    *uuid.UUID             `json:"secretTarget,omitempty"`    // target player for secret question
	ForAllResults   []ForAllAnswerResult   `json:"forAllResults,omitempty"`   // results for forAll question
}

// StakeInfo contains information for stake betting
type StakeInfo struct {
	MinBet     int  `json:"minBet"`     // Minimum bet (question price)
	MaxBet     int  `json:"maxBet"`     // Maximum bet (player's score or all-in)
	CurrentBet int  `json:"currentBet"` // Currently placed bet (0 if not placed yet)
	IsAllIn    bool `json:"isAllIn"`    // True if player bet all their points
}

// ForAllAnswerResult contains the result for a single player in forAll question
type ForAllAnswerResult struct {
	UserID     uuid.UUID `json:"userId"`
	Username   string    `json:"username"`
	Answer     string    `json:"answer"`
	IsCorrect  bool      `json:"isCorrect"`
	ScoreDelta int       `json:"scoreDelta"`
}

// ThemeState represents a theme with questions availability (DTO)
// All fields are required
type ThemeState struct {
	Name      string          `json:"name" binding:"required"`
	Questions []QuestionState `json:"questions" binding:"required"`
}

// QuestionState represents question availability in UI (DTO)
// Required: id, price, available, type
// Optional: text, mediaType, mediaUrl, mediaDurationMs (only when question is shown), answer (only for host)
type QuestionState struct {
	ID              string `json:"id" binding:"required"`
	Price           int    `json:"price" binding:"required"`
	Available       bool   `json:"available" binding:"required"`
	Type            string `json:"type" binding:"required"`   // normal, secret, stake, forAll
	Text            string `json:"text,omitempty"`            // only when shown
	MediaType       string `json:"mediaType,omitempty"`       // only when shown
	MediaURL        string `json:"mediaUrl,omitempty"`        // only when shown (image/audio/video URL)
	MediaDurationMs int    `json:"mediaDurationMs,omitempty"` // only for audio/video (duration in ms)
	Answer          string `json:"answer,omitempty"`          // only for host (correct answer)
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
// Required: user_id, username, role
// Optional: avatar_url
type PlayerInfo struct {
	UserID    uuid.UUID `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440002"`
	Username  string    `json:"username" binding:"required" example:"player1"`
	AvatarURL string    `json:"avatar_url,omitempty" example:"https://example.com/avatar.png"`
	Role      string    `json:"role" binding:"required" example:"player"`
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
