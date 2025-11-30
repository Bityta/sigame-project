package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

/**
 * Ответ на обновление настроек комнаты
 * Согласно README возвращает только объект settings
 */
data class UpdateRoomSettingsResponse(
    @JsonProperty("settings")
    val settings: RoomSettingsDto
)







