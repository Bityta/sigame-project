package redis

import "fmt"

func ErrMarshalGameState(err error) error {
	return fmt.Errorf("failed to marshal game state: %w", err)
}

func ErrSaveGameState(err error) error {
	return fmt.Errorf("failed to save game state: %w", err)
}

func ErrLoadGameState(err error) error {
	return fmt.Errorf("failed to load game state: %w", err)
}

func ErrUnmarshalGameState(err error) error {
	return fmt.Errorf("failed to unmarshal game state: %w", err)
}

func ErrGetAllScores(err error) error {
	return fmt.Errorf("failed to get all scores: %w", err)
}

func ErrParseScore(userID string, err error) error {
	return fmt.Errorf("failed to parse score for user %s: %w", userID, err)
}

func ErrGetActivePlayers(err error) error {
	return fmt.Errorf("failed to get active players: %w", err)
}

func ErrConvertFieldToString(field string, err error) error {
	return fmt.Errorf("failed to convert field %s to string: %w", field, err)
}

func ErrSetMetadataField(field string, err error) error {
	return fmt.Errorf("failed to set metadata field %s: %w", field, err)
}

func ErrMarshalValue(err error) error {
	return fmt.Errorf("failed to marshal value: %w", err)
}

func ErrGetActiveGames(err error) error {
	return fmt.Errorf("failed to get active games: %w", err)
}

func ErrParseUUID(uuidStr string, err error) error {
	return fmt.Errorf("failed to parse UUID %s: %w", uuidStr, err)
}

func ErrMarshalPack(err error) error {
	return fmt.Errorf("failed to marshal pack: %w", err)
}

func ErrCachePack(err error) error {
	return fmt.Errorf("failed to cache pack: %w", err)
}

var ErrPackNotFound = fmt.Errorf("pack not found")

func ErrGetCachedPack(err error) error {
	return fmt.Errorf("failed to get cached pack: %w", err)
}

func ErrUnmarshalPack(err error) error {
	return fmt.Errorf("failed to unmarshal pack: %w", err)
}

func ErrKeyNotFound(key string) error {
	return fmt.Errorf("key not found: %s", key)
}

func ErrGetValue(err error) error {
	return fmt.Errorf("failed to get value: %w", err)
}

