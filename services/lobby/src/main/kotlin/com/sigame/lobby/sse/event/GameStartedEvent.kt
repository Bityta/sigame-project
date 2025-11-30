package com.sigame.lobby.sse.event

import java.util.UUID

data class GameStartedEvent(
    val eventRoomId: UUID,
    val gameId: UUID,
    val websocketUrl: String
) : RoomEvent("game_started", eventRoomId)

