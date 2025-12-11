package com.sigame.lobby.sse.event

import com.sigame.lobby.domain.dto.RoomSettingsDto
import java.util.UUID

data class SettingsUpdatedEvent(
    val eventRoomId: UUID,
    val settings: RoomSettingsDto
) : RoomEvent("settings_updated", eventRoomId)
