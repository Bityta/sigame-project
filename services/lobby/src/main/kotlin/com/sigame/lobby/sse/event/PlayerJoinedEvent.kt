package com.sigame.lobby.sse.event

import java.util.UUID

data class PlayerJoinedEvent(
    val eventRoomId: UUID,
    val userId: UUID,
    val username: String,
    val currentPlayers: Int
) : RoomEvent("player_joined", eventRoomId)

