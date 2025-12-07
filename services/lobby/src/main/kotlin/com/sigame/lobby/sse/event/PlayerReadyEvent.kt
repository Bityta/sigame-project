package com.sigame.lobby.sse.event

import java.util.UUID

class PlayerReadyEvent(
    roomId: UUID,
    val userId: UUID,
    val username: String,
    val isReady: Boolean,
    val allPlayersReady: Boolean,
    val readyCount: Int,
    val totalCount: Int
) : RoomEvent("player_ready", roomId)

