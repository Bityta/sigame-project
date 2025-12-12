package postgres

import "fmt"

func ErrCreateGameSession(err error) error {
	return fmt.Errorf("failed to create game session: %w", err)
}

func ErrUpdateGameSession(err error) error {
	return fmt.Errorf("failed to update game session: %w", err)
}

func ErrGetGameSession(err error) error {
	return fmt.Errorf("failed to get game session: %w", err)
}

func ErrBeginTransaction(err error) error {
	return fmt.Errorf("failed to begin transaction: %w", err)
}

func ErrUpdatePlayerScore(err error) error {
	return fmt.Errorf("failed to update player score: %w", err)
}

func ErrCommitTransaction(err error) error {
	return fmt.Errorf("failed to commit transaction: %w", err)
}

func ErrGetGames(err error) error {
	return fmt.Errorf("failed to get games: %w", err)
}

func ErrGetActiveGameForUser(err error) error {
	return fmt.Errorf("failed to get active game for user: %w", err)
}

func ErrLogEvent(err error) error {
	return fmt.Errorf("failed to log event: %w", err)
}

func ErrPrepareStatement(err error) error {
	return fmt.Errorf("failed to prepare statement: %w", err)
}

func ErrInsertEvent(err error) error {
	return fmt.Errorf("failed to insert event: %w", err)
}

func ErrGetEvents(err error) error {
	return fmt.Errorf("failed to get events: %w", err)
}

func ErrGetEventCount(err error) error {
	return fmt.Errorf("failed to get event count: %w", err)
}

