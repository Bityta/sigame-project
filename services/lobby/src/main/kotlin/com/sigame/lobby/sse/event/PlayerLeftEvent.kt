package com.sigame.lobby.sse.event

import java.util.UUID

data class PlayerLeftEvent(
    val eventRoomId: UUID,
    val userId: UUID,
    val username: String,
    val reason: String,
    val currentPlayers: Int
) : RoomEvent("player_left", eventRoomId)

