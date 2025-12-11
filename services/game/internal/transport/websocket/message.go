package websocket

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain"
)

// MessageType represents the type of WebSocket message
type MessageType string

// WebSocket message type constants (client to server and server to client)
const (
	// Client -> Server messages
	// MessageTypeReady indicates player is ready
	MessageTypeReady MessageType = "READY"
	// MessageTypeSelectQuestion indicates question selection
	MessageTypeSelectQuestion MessageType = "SELECT_QUESTION"
	// MessageTypePressButton indicates button press
	MessageTypePressButton MessageType = "PRESS_BUTTON"
	// MessageTypeSubmitAnswer indicates answer submission
	MessageTypeSubmitAnswer MessageType = "SUBMIT_ANSWER"
	// MessageTypeJudgeAnswer indicates answer judging
	MessageTypeJudgeAnswer MessageType = "JUDGE_ANSWER"
	// MessageTypeMediaLoadProgress indicates media loading progress
	MessageTypeMediaLoadProgress MessageType = "MEDIA_LOAD_PROGRESS"
	// MessageTypeMediaLoadComplete indicates all media loaded
	MessageTypeMediaLoadComplete MessageType = "MEDIA_LOAD_COMPLETE"

	// Server -> Client messages
	// MessageTypeStateUpdate indicates game state update
	MessageTypeStateUpdate MessageType = "STATE_UPDATE"
	// MessageTypeQuestionSelected indicates question was selected
	MessageTypeQuestionSelected MessageType = "QUESTION_SELECTED"
	// MessageTypeButtonPressed indicates button was pressed
	MessageTypeButtonPressed MessageType = "BUTTON_PRESSED"
	// MessageTypeAnswerResult indicates answer result
	MessageTypeAnswerResult MessageType = "ANSWER_RESULT"
	// MessageTypeRoundComplete indicates round completion
	MessageTypeRoundComplete MessageType = "ROUND_COMPLETE"
	// MessageTypeGameComplete indicates game completion
	MessageTypeGameComplete MessageType = "GAME_COMPLETE"
	// MessageTypeError indicates an error occurred
	MessageTypeError MessageType = "ERROR"
	// MessageTypeRoundMediaManifest contains all media for the round
	MessageTypeRoundMediaManifest MessageType = "ROUND_MEDIA_MANIFEST"
	// MessageTypeStartMedia commands synchronized media playback
	MessageTypeStartMedia MessageType = "START_MEDIA"
)

// ClientMessage represents a message from client to server
type ClientMessage struct {
	Type    MessageType            `json:"type"`
	UserID  uuid.UUID              `json:"user_id"`
	GameID  uuid.UUID              `json:"game_id"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// ServerMessage represents a message from server to client
type ServerMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

// SelectQuestionPayload represents the payload for selecting a question
type SelectQuestionPayload struct {
	ThemeID    string `json:"theme_id"`
	QuestionID string `json:"question_id"`
}

// SubmitAnswerPayload represents the payload for submitting an answer
type SubmitAnswerPayload struct {
	Answer string `json:"answer"`
}

// JudgeAnswerPayload represents the payload for judging an answer
type JudgeAnswerPayload struct {
	UserID  uuid.UUID `json:"user_id"`
	Correct bool      `json:"correct"`
}

// ButtonPressedPayload represents the payload when a button is pressed
type ButtonPressedPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	LatencyMS int64     `json:"latency_ms"`
}

// AnswerResultPayload represents the payload for answer result
type AnswerResultPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Correct   bool      `json:"correct"`
	Answer    string    `json:"answer,omitempty"`
	Score     int       `json:"score"`
	ScoreDelta int      `json:"score_delta"`
}

// RoundCompletePayload represents the payload for round completion
type RoundCompletePayload struct {
	RoundNumber int                  `json:"round_number"`
	Scores      []domain.PlayerScore `json:"scores"`
	NextRound   *int                 `json:"next_round,omitempty"`
}

// GameCompletePayload represents the payload for game completion
type GameCompletePayload struct {
	Winners []domain.PlayerScore `json:"winners"`
	Scores  []domain.PlayerScore `json:"scores"`
}

// ErrorPayload represents an error message
type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// MediaItem represents a single media file in the manifest
type MediaItem struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"` // "image", "audio", "video"
	URL          string       `json:"url"`
	Size         int64        `json:"size"` // bytes
	QuestionRef  QuestionRef  `json:"question_ref"`
}

// QuestionRef identifies which question the media belongs to
type QuestionRef struct {
	ThemeIndex    int `json:"theme"`
	QuestionPrice int `json:"price"`
}

// RoundMediaManifestPayload contains all media for preloading
type RoundMediaManifestPayload struct {
	Round      int         `json:"round"`
	Media      []MediaItem `json:"media"`
	TotalSize  int64       `json:"total_size"`
	TotalCount int         `json:"total_count"`
}

// MediaLoadProgressPayload represents client's media loading progress
type MediaLoadProgressPayload struct {
	Loaded      int   `json:"loaded"`
	Total       int   `json:"total"`
	BytesLoaded int64 `json:"bytes_loaded"`
	Percent     int   `json:"percent"`
}

// MediaLoadCompletePayload indicates all media loaded
type MediaLoadCompletePayload struct {
	Round       int `json:"round"`
	LoadedCount int `json:"loaded_count"`
}

// StartMediaPayload commands synchronized media playback
type StartMediaPayload struct {
	MediaID    string `json:"media_id"`
	MediaType  string `json:"media_type"` // "image", "audio", "video"
	URL        string `json:"url"`
	StartAt    int64  `json:"start_at"`     // Unix timestamp in ms when to start
	DurationMS int64  `json:"duration_ms"`  // Media duration for sync
}

// NewClientMessage creates a ClientMessage from JSON bytes
func NewClientMessage(data []byte) (*ClientMessage, error) {
	var msg ClientMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// NewServerMessage creates a new server message
func NewServerMessage(msgType MessageType, payload interface{}) *ServerMessage {
	return &ServerMessage{
		Type:    msgType,
		Payload: payload,
	}
}

// ToJSON converts ServerMessage to JSON bytes
func (m *ServerMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// NewStateUpdateMessage creates a state update message
func NewStateUpdateMessage(state *domain.GameState) *ServerMessage {
	return NewServerMessage(MessageTypeStateUpdate, state)
}

// NewButtonPressedMessage creates a button pressed message
func NewButtonPressedMessage(userID uuid.UUID, username string, latencyMS int64) *ServerMessage {
	return NewServerMessage(MessageTypeButtonPressed, ButtonPressedPayload{
		UserID:    userID,
		Username:  username,
		LatencyMS: latencyMS,
	})
}

// NewAnswerResultMessage creates an answer result message
func NewAnswerResultMessage(userID uuid.UUID, username string, correct bool, answer string, score, scoreDelta int) *ServerMessage {
	return NewServerMessage(MessageTypeAnswerResult, AnswerResultPayload{
		UserID:     userID,
		Username:   username,
		Correct:    correct,
		Answer:     answer,
		Score:      score,
		ScoreDelta: scoreDelta,
	})
}

// NewErrorMessage creates an error message
func NewErrorMessage(message, code string) *ServerMessage {
	return NewServerMessage(MessageTypeError, ErrorPayload{
		Message: message,
		Code:    code,
	})
}

// NewRoundMediaManifestMessage creates a media manifest message
func NewRoundMediaManifestMessage(round int, media []MediaItem, totalSize int64) *ServerMessage {
	return NewServerMessage(MessageTypeRoundMediaManifest, RoundMediaManifestPayload{
		Round:      round,
		Media:      media,
		TotalSize:  totalSize,
		TotalCount: len(media),
	})
}

// NewStartMediaMessage creates a start media message for synchronized playback
func NewStartMediaMessage(mediaID, mediaType, url string, startAt, durationMS int64) *ServerMessage {
	return NewServerMessage(MessageTypeStartMedia, StartMediaPayload{
		MediaID:    mediaID,
		MediaType:  mediaType,
		URL:        url,
		StartAt:    startAt,
		DurationMS: durationMS,
	})
}

