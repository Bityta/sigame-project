package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

/**
 * Ответ на запрос старта игры
 * Согласно README: gameId и websocketUrl
 */
data class StartGameResponse(
    @JsonProperty("gameId")
    val gameId: String,
    
    @JsonProperty("websocketUrl")
    val websocketUrl: String
)

