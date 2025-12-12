package game

import "fmt"

var (
	ErrClientDoesNotImplementSend = fmt.Errorf("client does not implement Send method")
	ErrInvalidManifestType        = fmt.Errorf("invalid manifest type")
	ErrRoundNotFound              = fmt.Errorf("round not found")
	ErrThemeNotFound              = fmt.Errorf("theme not found")
	ErrMediaItemNotFound          = fmt.Errorf("media item not found")
)

func ErrSerializeState(err error) error {
	return fmt.Errorf("failed to serialize state: %w", err)
}

func ErrSerializeStateForClient(err error) error {
	return fmt.Errorf("failed to serialize state for client: %w", err)
}

func ErrSerializeMediaManifest(err error) error {
	return fmt.Errorf("failed to serialize media manifest: %w", err)
}

func ErrSerializeStartMediaMessage(err error) error {
	return fmt.Errorf("failed to serialize start media message: %w", err)
}

