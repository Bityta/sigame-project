package domain

import "errors"

var (
	// Game errors
	ErrGameNotFound      = errors.New("game not found")
	ErrGameAlreadyExists = errors.New("game already exists")
	ErrGameFinished      = errors.New("game already finished")
	ErrGameNotStarted    = errors.New("game not started")
	ErrInvalidGameState  = errors.New("invalid game state")

	// Player errors
	ErrPlayerNotFound    = errors.New("player not found")
	ErrPlayerNotActive   = errors.New("player not active")
	ErrNotPlayerTurn     = errors.New("not player's turn")
	ErrPlayerNotReady    = errors.New("player not ready")
	ErrInsufficientPlayers = errors.New("insufficient players to start game")

	// Question errors
	ErrQuestionNotFound    = errors.New("question not found")
	ErrQuestionAlreadyUsed = errors.New("question already used")
	ErrInvalidAnswer       = errors.New("invalid answer")

	// Round errors
	ErrRoundNotFound  = errors.New("round not found")
	ErrRoundComplete  = errors.New("round already complete")
	ErrInvalidRound   = errors.New("invalid round number")

	// Button press errors
	ErrButtonAlreadyPressed = errors.New("button already pressed")
	ErrButtonPressTimeout   = errors.New("button press timeout")

	// Pack errors
	ErrPackNotFound  = errors.New("pack not found")
	ErrPackInvalid   = errors.New("pack is invalid")

	// WebSocket errors
	ErrWSConnectionClosed = errors.New("websocket connection closed")
	ErrWSInvalidMessage   = errors.New("invalid websocket message")

	// Timer errors
	ErrTimerExpired = errors.New("timer expired")
)

