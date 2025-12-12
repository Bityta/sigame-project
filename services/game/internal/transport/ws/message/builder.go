package message

import (
	"encoding/json"

	"github.com/google/uuid"
	domainGame "github.com/sigame/game/internal/domain/game"
)

func NewClientMessage(data []byte) (*ClientMessage, error) {
	var msg ClientMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func NewServerMessage(msgType MessageType, payload interface{}) *ServerMessage {
	return &ServerMessage{
		Type:    msgType,
		Payload: payload,
	}
}

func (m *ServerMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func NewStateUpdateMessage(state *domainGame.State) *ServerMessage {
	return NewServerMessage(MessageTypeStateUpdate, state)
}

func NewPingMessage(serverTime int64) *ServerMessage {
	return NewServerMessage(MessageTypePing, PingPayload{
		ServerTime: serverTime,
	})
}

func NewButtonPressedMessage(winnerID uuid.UUID, winnerName string, reactionTimeMS int64, allPresses []PressInfo) *ServerMessage {
	return NewServerMessage(MessageTypeButtonPressed, ButtonPressedPayload{
		WinnerID:       winnerID,
		WinnerName:     winnerName,
		ReactionTimeMS: reactionTimeMS,
		AllPresses:     allPresses,
	})
}

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

func NewErrorMessage(message, code string) *ServerMessage {
	return NewServerMessage(MessageTypeError, ErrorPayload{
		Message: message,
		Code:    code,
	})
}

func NewRoundMediaManifestMessage(round int, media []MediaItem, totalSize int64) *ServerMessage {
	return NewServerMessage(MessageTypeRoundMediaManifest, RoundMediaManifestPayload{
		Round:      round,
		Media:      media,
		TotalSize:  totalSize,
		TotalCount: len(media),
	})
}

func NewStartMediaMessage(mediaID, mediaType, url string, startAt, durationMS int64) *ServerMessage {
	return NewServerMessage(MessageTypeStartMedia, StartMediaPayload{
		MediaID:    mediaID,
		MediaType:  mediaType,
		URL:        url,
		StartAt:    startAt,
		DurationMS: durationMS,
	})
}

func NewSecretTransferredMessage(fromUserID uuid.UUID, fromUsername string, toUserID uuid.UUID, toUsername string) *ServerMessage {
	return NewServerMessage(MessageTypeSecretTransferred, SecretTransferredPayload{
		FromUserID:   fromUserID,
		FromUsername: fromUsername,
		ToUserID:     toUserID,
		ToUsername:   toUsername,
	})
}

func NewStakePlacedMessage(userID uuid.UUID, username string, amount int, allIn bool) *ServerMessage {
	return NewServerMessage(MessageTypeStakePlaced, StakePlacedPayload{
		UserID:   userID,
		Username: username,
		Amount:   amount,
		AllIn:    allIn,
	})
}

func NewForAllResultsMessage(correctAnswer string, results []domainGame.ForAllResult) *ServerMessage {
	return NewServerMessage(MessageTypeForAllResults, ForAllResultsPayload{
		CorrectAnswer: correctAnswer,
		Results:       results,
	})
}

