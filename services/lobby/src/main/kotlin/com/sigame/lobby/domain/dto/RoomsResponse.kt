package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

data class RoomsResponse(
    @JsonProperty("rooms")
    val rooms: List<RoomDto>
)

