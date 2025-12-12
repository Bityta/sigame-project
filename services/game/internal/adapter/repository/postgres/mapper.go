package postgres

import (
	"github.com/sigame/game/internal/domain/event"
	domainGame "github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
)

func toDomainGame(dbGame interface{}) *domainGame.Game {
	return nil
}

func toDBGame(game *domainGame.Game) interface{} {
	return nil
}

func toDomainPlayer(dbPlayer interface{}) *player.Player {
	return nil
}

func toDBPlayer(p *player.Player) interface{} {
	return nil
}

func toDomainEvent(dbEvent interface{}) *event.Event {
	return nil
}

func toDBEvent(e *event.Event) interface{} {
	return nil
}

func toDomainPack(dbPack interface{}) *pack.Pack {
	return nil
}

