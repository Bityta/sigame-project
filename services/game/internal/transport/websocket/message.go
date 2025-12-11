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
	// MessageTypeSelectQuestion indicates question selection
	MessageTypeSelectQuestion MessageType = "SELECT_QUESTION"
	// MessageTypePressButton indicates button press
	MessageTypePressButton MessageType = "PRESS_BUTTON"
	// MessageTypeSubmitAnswer indicates answer submission
	MessageTypeSubmitAnswer MessageType = "SUBMIT_ANSWER"
	// MessageTypeJudgeAnswer indicates answer judging
	MessageTypeJudgeAnswer MessageType = "JUDGE_ANSWER"
	// MessageTypePong is client response to PING for RTT measurement
	MessageTypePong MessageType = "PONG"
	// MessageTypeMediaLoadProgress indicates media loading progress
	MessageTypeMediaLoadProgress MessageType = "MEDIA_LOAD_PROGRESS"
	// MessageTypeMediaLoadComplete indicates all media loaded
	MessageTypeMediaLoadComplete MessageType = "MEDIA_LOAD_COMPLETE"
	// MessageTypeTransferSecret - host transfers secret question to player
	MessageTypeTransferSecret MessageType = "TRANSFER_SECRET"
	// MessageTypePlaceStake - player places a stake for stake question
	MessageTypePlaceStake MessageType = "PLACE_STAKE"
	// MessageTypeSubmitForAllAnswer - player submits answer for forAll question
	MessageTypeSubmitForAllAnswer MessageType = "SUBMIT_FOR_ALL_ANSWER"

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
	// MessageTypePing is server ping for RTT measurement
	MessageTypePing MessageType = "PING"
	// MessageTypeRoundMediaManifest contains all media for the round
	MessageTypeRoundMediaManifest MessageType = "ROUND_MEDIA_MANIFEST"
	// MessageTypeStartMedia commands synchronized media playback
	MessageTypeStartMedia MessageType = "START_MEDIA"
	// MessageTypeSecretTransferred - notification that secret question was transferred
	MessageTypeSecretTransferred MessageType = "SECRET_TRANSFERRED"
	// MessageTypeStakePlaced - notification that stake was placed
	MessageTypeStakePlaced MessageType = "STAKE_PLACED"
	// MessageTypeForAllResults - results of forAll question
	MessageTypeForAllResults MessageType = "FOR_ALL_RESULTS"
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

// PingPayload represents the payload for PING message (server to client)
type PingPayload struct {
	ServerTime int64 `json:"server_time"`
}

// PongPayload represents the payload for PONG message (client to server)
type PongPayload struct {
	ServerTime int64 `json:"server_time"`
	ClientTime int64 `json:"client_time"`
}

// PressInfo represents information about a single button press
type PressInfo struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	TimeMS   int64     `json:"time_ms"` // Adjusted reaction time in milliseconds
}

// ButtonPressedPayload represents the payload when a button is pressed
type ButtonPressedPayload struct {
	WinnerID       uuid.UUID   `json:"winner_id"`
	WinnerName     string      `json:"winner_name"`
	ReactionTimeMS int64       `json:"reaction_time_ms"` // Winner's adjusted reaction time
	AllPresses     []PressInfo `json:"all_presses"`      // All button presses sorted by adjusted time
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

// TransferSecretPayload represents the payload for transferring secret question
type TransferSecretPayload struct {
	TargetUserID uuid.UUID `json:"target_user_id"`
}

// PlaceStakePayload represents the payload for placing a stake
type PlaceStakePayload struct {
	Amount int  `json:"amount"`
	AllIn  bool `json:"all_in"` // If true, bet all points
}

// SubmitForAllAnswerPayload represents the payload for forAll answer
type SubmitForAllAnswerPayload struct {
	Answer string `json:"answer"`
}

// SecretTransferredPayload notification that secret was transferred
type SecretTransferredPayload struct {
	FromUserID   uuid.UUID `json:"from_user_id"`
	FromUsername string    `json:"from_username"`
	ToUserID     uuid.UUID `json:"to_user_id"`
	ToUsername   string    `json:"to_username"`
}

// StakePlacedPayload notification that stake was placed
type StakePlacedPayload struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Amount   int       `json:"amount"`
	AllIn    bool      `json:"all_in"`
}

// ForAllResultsPayload contains results for all players
type ForAllResultsPayload struct {
	CorrectAnswer string                      `json:"correct_answer"`
	Results       []domain.ForAllAnswerResult `json:"results"`
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

// NewPingMessage creates a PING message for RTT measurement
func NewPingMessage(serverTime int64) *ServerMessage {
	return NewServerMessage(MessageTypePing, PingPayload{
		ServerTime: serverTime,
	})
}

// NewButtonPressedMessage creates a button pressed message with all press info
func NewButtonPressedMessage(winnerID uuid.UUID, winnerName string, reactionTimeMS int64, allPresses []PressInfo) *ServerMessage {
	return NewServerMessage(MessageTypeButtonPressed, ButtonPressedPayload{
		WinnerID:       winnerID,
		WinnerName:     winnerName,
		ReactionTimeMS: reactionTimeMS,
		AllPresses:     allPresses,
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

// NewSecretTransferredMessage creates a secret transferred notification
func NewSecretTransferredMessage(fromUserID uuid.UUID, fromUsername string, toUserID uuid.UUID, toUsername string) *ServerMessage {
	return NewServerMessage(MessageTypeSecretTransferred, SecretTransferredPayload{
		FromUserID:   fromUserID,
		FromUsername: fromUsername,
		ToUserID:     toUserID,
		ToUsername:   toUsername,
	})
}

// NewStakePlacedMessage creates a stake placed notification
func NewStakePlacedMessage(userID uuid.UUID, username string, amount int, allIn bool) *ServerMessage {
	return NewServerMessage(MessageTypeStakePlaced, StakePlacedPayload{
		UserID:   userID,
		Username: username,
		Amount:   amount,
		AllIn:    allIn,
	})
}

// NewForAllResultsMessage creates a forAll results message
func NewForAllResultsMessage(correctAnswer string, results []domain.ForAllAnswerResult) *ServerMessage {
	return NewServerMessage(MessageTypeForAllResults, ForAllResultsPayload{
		CorrectAnswer: correctAnswer,
		Results:       results,
	})
}
