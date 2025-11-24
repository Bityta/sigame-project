package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

/**
 * Ответ со списком комнат с пагинацией
 * Все поля обязательны
 */
data class RoomListResponse(
    @JsonProperty("rooms")
    val rooms: List<RoomDto>,
    
    @JsonProperty("total")
    val total: Long,
    
    @JsonProperty("page")
    val page: Int,
    
    @JsonProperty("size")
    val size: Int
)

