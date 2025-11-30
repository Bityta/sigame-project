package com.sigame.lobby.sse.event

import java.util.UUID

data class SettingsUpdatedEvent(
    val eventRoomId: UUID
) : RoomEvent("settings_updated", eventRoomId)

