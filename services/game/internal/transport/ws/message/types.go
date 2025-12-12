package message

import (
	"github.com/google/uuid"
	domainGame "sigame/game/internal/domain/game"
	"sigame/game/internal/domain/player"
)

type MessageType string

const (
	MessageTypeSelectQuestion MessageType = "SELECT_QUESTION"
	MessageTypePressButton MessageType = "PRESS_BUTTON"
	MessageTypeSubmitAnswer MessageType = "SUBMIT_ANSWER"
	MessageTypeJudgeAnswer MessageType = "JUDGE_ANSWER"
	MessageTypePong MessageType = "PONG"
	MessageTypeMediaLoadProgress MessageType = "MEDIA_LOAD_PROGRESS"
	MessageTypeMediaLoadComplete MessageType = "MEDIA_LOAD_COMPLETE"
	MessageTypeTransferSecret MessageType = "TRANSFER_SECRET"
	MessageTypePlaceStake MessageType = "PLACE_STAKE"
	MessageTypeSubmitForAllAnswer MessageType = "SUBMIT_FOR_ALL_ANSWER"

	MessageTypeStateUpdate MessageType = "STATE_UPDATE"
	MessageTypeQuestionSelected MessageType = "QUESTION_SELECTED"
	MessageTypeButtonPressed MessageType = "BUTTON_PRESSED"
	MessageTypeAnswerResult MessageType = "ANSWER_RESULT"
	MessageTypeRoundComplete MessageType = "ROUND_COMPLETE"
	MessageTypeGameComplete MessageType = "GAME_COMPLETE"
	MessageTypeError MessageType = "ERROR"
	MessageTypePing MessageType = "PING"
	MessageTypeRoundMediaManifest MessageType = "ROUND_MEDIA_MANIFEST"
	MessageTypeStartMedia MessageType = "START_MEDIA"
	MessageTypeSecretTransferred MessageType = "SECRET_TRANSFERRED"
	MessageTypeStakePlaced MessageType = "STAKE_PLACED"
	MessageTypeForAllResults MessageType = "FOR_ALL_RESULTS"
)

type ClientMessage struct {
	Type    MessageType            `json:"type"`
	UserID  uuid.UUID              `json:"user_id"`
	GameID  uuid.UUID              `json:"game_id"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func (m *ClientMessage) GetType() string {
	return string(m.Type)
}

func (m *ClientMessage) GetPayload() map[string]interface{} {
	return m.Payload
}

type ServerMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

type SelectQuestionPayload struct {
	ThemeID    string `json:"theme_id"`
	QuestionID string `json:"question_id"`
}

type SubmitAnswerPayload struct {
	Answer string `json:"answer"`
}

type JudgeAnswerPayload struct {
	UserID  uuid.UUID `json:"user_id"`
	Correct bool      `json:"correct"`
}

type PingPayload struct {
	ServerTime int64 `json:"server_time"`
}

type PongPayload struct {
	ServerTime int64 `json:"server_time"`
	ClientTime int64 `json:"client_time"`
}

type PressInfo struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	TimeMS   int64     `json:"time_ms"`
}

type ButtonPressedPayload struct {
	WinnerID       uuid.UUID   `json:"winner_id"`
	WinnerName     string      `json:"winner_name"`
	ReactionTimeMS int64       `json:"reaction_time_ms"`
	AllPresses     []PressInfo `json:"all_presses"`
}

type AnswerResultPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Correct   bool      `json:"correct"`
	Answer    string    `json:"answer,omitempty"`
	Score     int       `json:"score"`
	ScoreDelta int      `json:"score_delta"`
}

type RoundCompletePayload struct {
	RoundNumber int            `json:"round_number"`
	Scores      []player.Score `json:"scores"`
	NextRound   *int           `json:"next_round,omitempty"`
}

type GameCompletePayload struct {
	Winners []player.Score `json:"winners"`
	Scores  []player.Score `json:"scores"`
}

type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type MediaItem struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	URL          string       `json:"url"`
	Size         int64        `json:"size"`
	QuestionRef  QuestionRef  `json:"question_ref"`
}

type QuestionRef struct {
	ThemeIndex    int `json:"theme"`
	QuestionPrice int `json:"price"`
}

type RoundMediaManifestPayload struct {
	Round      int         `json:"round"`
	Media      []MediaItem `json:"media"`
	TotalSize  int64       `json:"total_size"`
	TotalCount int         `json:"total_count"`
}

type MediaLoadProgressPayload struct {
	Loaded      int   `json:"loaded"`
	Total       int   `json:"total"`
	BytesLoaded int64 `json:"bytes_loaded"`
	Percent     int   `json:"percent"`
}

type MediaLoadCompletePayload struct {
	Round       int `json:"round"`
	LoadedCount int `json:"loaded_count"`
}

type StartMediaPayload struct {
	MediaID    string `json:"media_id"`
	MediaType  string `json:"media_type"`
	URL        string `json:"url"`
	StartAt    int64  `json:"start_at"`
	DurationMS int64  `json:"duration_ms"`
}

type TransferSecretPayload struct {
	TargetUserID uuid.UUID `json:"target_user_id"`
}

type PlaceStakePayload struct {
	Amount int  `json:"amount"`
	AllIn  bool `json:"all_in"`
}

type SubmitForAllAnswerPayload struct {
	Answer string `json:"answer"`
}

type SecretTransferredPayload struct {
	FromUserID   uuid.UUID `json:"from_user_id"`
	FromUsername string    `json:"from_username"`
	ToUserID     uuid.UUID `json:"to_user_id"`
	ToUsername   string    `json:"to_username"`
}

type StakePlacedPayload struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Amount   int       `json:"amount"`
	AllIn    bool      `json:"all_in"`
}

type ForAllResultsPayload struct {
	CorrectAnswer string                   `json:"correct_answer"`
	Results       []domainGame.ForAllResult `json:"results"`
}

