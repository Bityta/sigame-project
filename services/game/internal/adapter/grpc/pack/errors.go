package pack

import "fmt"

func ErrPackServiceError(errMsg string) error {
	return fmt.Errorf("pack service returned error: %s", errMsg)
}

func ErrCreateRequest(err error) error {
	return fmt.Errorf("failed to create request: %w", err)
}

func ErrExecuteRequest(err error) error {
	return fmt.Errorf("failed to execute request: %w", err)
}

func ErrReadResponse(err error) error {
	return fmt.Errorf("failed to read response body: %w", err)
}

