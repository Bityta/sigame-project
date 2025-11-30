package com.sigame.lobby.sse.event

import java.util.UUID

data class RoomClosedEvent(
    val eventRoomId: UUID,
    val reason: String
) : RoomEvent("room_closed", eventRoomId)

