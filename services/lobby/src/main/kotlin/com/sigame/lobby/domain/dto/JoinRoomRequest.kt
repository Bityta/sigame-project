package com.sigame.lobby.domain.dto

/**
 * Запрос на присоединение к комнате
 */
data class JoinRoomRequest(
    val password: String? = null,
    val role: String = "player"
)

