package domain

import "errors"

var (
	// ErrGameNotFound indicates a game was not found
	ErrGameNotFound = errors.New("game not found")
	// ErrGameAlreadyExists indicates a game already exists
	ErrGameAlreadyExists = errors.New("game already exists")
	// ErrGameFinished indicates game has finished
	ErrGameFinished = errors.New("game already finished")
	// ErrGameNotStarted indicates game has not started yet
	ErrGameNotStarted = errors.New("game not started")
	// ErrInvalidGameState indicates invalid game state
	ErrInvalidGameState = errors.New("invalid game state")

	// ErrPlayerNotFound indicates player was not found
	ErrPlayerNotFound = errors.New("player not found")
	// ErrPlayerNotActive indicates player is not active
	ErrPlayerNotActive = errors.New("player not active")
	// ErrNotPlayerTurn indicates it's not the player's turn
	ErrNotPlayerTurn = errors.New("not player's turn")
	// ErrPlayerNotReady indicates player is not ready
	ErrPlayerNotReady = errors.New("player not ready")
	// ErrInsufficientPlayers indicates not enough players
	ErrInsufficientPlayers = errors.New("insufficient players to start game")

	// ErrQuestionNotFound indicates question was not found
	ErrQuestionNotFound = errors.New("question not found")
	// ErrQuestionAlreadyUsed indicates question was already used
	ErrQuestionAlreadyUsed = errors.New("question already used")
	// ErrInvalidAnswer indicates answer is invalid
	ErrInvalidAnswer = errors.New("invalid answer")

	// ErrRoundNotFound indicates round was not found
	ErrRoundNotFound = errors.New("round not found")
	// ErrRoundComplete indicates round is complete
	ErrRoundComplete = errors.New("round already complete")
	// ErrInvalidRound indicates invalid round number
	ErrInvalidRound = errors.New("invalid round number")

	// ErrButtonAlreadyPressed indicates button was already pressed
	ErrButtonAlreadyPressed = errors.New("button already pressed")
	// ErrButtonPressTimeout indicates button press timeout
	ErrButtonPressTimeout = errors.New("button press timeout")

	// ErrPackNotFound indicates pack was not found
	ErrPackNotFound = errors.New("pack not found")
	// ErrPackInvalid indicates pack is invalid
	ErrPackInvalid = errors.New("pack is invalid")

	// ErrWSConnectionClosed indicates websocket connection closed
	ErrWSConnectionClosed = errors.New("websocket connection closed")
	// ErrWSInvalidMessage indicates invalid websocket message
	ErrWSInvalidMessage = errors.New("invalid websocket message")

	// ErrTimerExpired indicates timer has expired
	ErrTimerExpired = errors.New("timer expired")
)

