package game

import "errors"

var (
	ErrGameNotFound        = errors.New("game not found")
	ErrInvalidGameStatus   = errors.New("invalid game status")
	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrPlayerNotFound      = errors.New("player not found")
	ErrInvalidRound        = errors.New("invalid round number")
	ErrHostNotFound        = errors.New("host not found")
	ErrInvalidSettings     = errors.New("invalid game settings")
)

