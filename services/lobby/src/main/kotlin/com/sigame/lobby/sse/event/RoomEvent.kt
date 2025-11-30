package com.sigame.lobby.sse.event

import java.time.Instant
import java.util.UUID

sealed class RoomEvent(
    val type: String,
    val roomId: UUID,
    val timestamp: Instant = Instant.now()
)

