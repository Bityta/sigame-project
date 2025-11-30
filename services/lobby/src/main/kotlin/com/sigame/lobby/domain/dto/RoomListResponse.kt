package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

data class RoomListResponse(
    @JsonProperty("rooms")
    val rooms: List<RoomDto>,
    
    @JsonProperty("page")
    val page: Int,
    
    @JsonProperty("size")
    val size: Int,
    
    @JsonProperty("totalElements")
    val totalElements: Long,
    
    @JsonProperty("totalPages")
    val totalPages: Int
)

