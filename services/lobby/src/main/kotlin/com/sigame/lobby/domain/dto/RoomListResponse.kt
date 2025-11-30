package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

/**
 * Ответ со списком комнат с пагинацией
 * Согласно README: rooms, page, size, totalElements, totalPages
 */
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

