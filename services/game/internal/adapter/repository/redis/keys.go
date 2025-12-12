package redis

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	keyPrefixPack      = "pack"
	keyPrefixGame      = "game"
	keyPrefixGames     = "games"
	keySuffixContent   = "content"
	keySuffixState     = "state"
	keySuffixScores    = "scores"
	keySuffixPlayers   = "players"
	keySuffixMeta      = "meta"
	keySuffixActive    = "active"
)

func packKey(packID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefixPack, packID.String(), keySuffixContent)
}

func gameStateKey(gameID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefixGame, gameID.String(), keySuffixState)
}

func gameScoresKey(gameID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefixGame, gameID.String(), keySuffixScores)
}

func gamePlayersKey(gameID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefixGame, gameID.String(), keySuffixPlayers)
}

func gameMetadataKey(gameID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefixGame, gameID.String(), keySuffixMeta)
}

func activeGamesKey() string {
	return fmt.Sprintf("%s:%s", keyPrefixGames, keySuffixActive)
}

