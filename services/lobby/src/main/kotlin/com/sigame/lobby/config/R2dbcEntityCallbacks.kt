package com.sigame.lobby.config

import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomSettings
import org.reactivestreams.Publisher
import org.springframework.data.r2dbc.mapping.event.AfterConvertCallback
import org.springframework.data.relational.core.sql.SqlIdentifier
import org.springframework.stereotype.Component
import reactor.core.publisher.Mono

@Component
class GameRoomAfterConvertCallback : AfterConvertCallback<GameRoom> {
    override fun onAfterConvert(entity: GameRoom, table: SqlIdentifier): Publisher<GameRoom> {
        return Mono.just(entity.markPersisted())
    }
}

@Component
class RoomSettingsAfterConvertCallback : AfterConvertCallback<RoomSettings> {
    override fun onAfterConvert(entity: RoomSettings, table: SqlIdentifier): Publisher<RoomSettings> {
        return Mono.just(entity.markPersisted())
    }
}

