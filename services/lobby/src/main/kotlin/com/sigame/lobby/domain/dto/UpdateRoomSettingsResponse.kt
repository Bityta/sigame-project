package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

data class UpdateRoomSettingsResponse(
    @JsonProperty("settings")
    val settings: RoomSettingsDto
)







