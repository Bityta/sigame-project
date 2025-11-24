package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

/**
 * Ответ на запрос старта игры
 * Все поля обязательны
 */
data class StartGameResponse(
    @JsonProperty("gameSessionId")
    val gameSessionId: String,
    
    @JsonProperty("wsUrl")
    val wsUrl: String
)

