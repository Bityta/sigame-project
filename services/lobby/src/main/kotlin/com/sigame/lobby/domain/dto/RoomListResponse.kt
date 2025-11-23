package com.sigame.lobby.domain.dto

/**
 * Ответ со списком комнат с пагинацией
 */
data class RoomListResponse(
    val rooms: List<RoomDto>,
    val total: Long,
    val page: Int,
    val size: Int
)

