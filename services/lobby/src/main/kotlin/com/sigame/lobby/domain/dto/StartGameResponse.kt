package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

data class StartGameResponse(
    @JsonProperty("gameId")
    val gameId: String,
    
    @JsonProperty("websocketUrl")
    val websocketUrl: String
)

