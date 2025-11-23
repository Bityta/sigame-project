package com.sigame.lobby.domain.dto

/**
 * Ответ на запрос старта игры
 */
data class StartGameResponse(
    val gameSessionId: String,
    val wsUrl: String
)

