package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

data class SetReadyResponse(
    @JsonProperty("isReady")
    val isReady: Boolean,

    @JsonProperty("allPlayersReady")
    val allPlayersReady: Boolean,

    @JsonProperty("readyCount")
    val readyCount: Int,

    @JsonProperty("totalCount")
    val totalCount: Int,

    @JsonProperty("gameStarted")
    val gameStarted: Boolean = false,

    @JsonProperty("gameId")
    val gameId: String? = null,

    @JsonProperty("websocketUrl")
    val websocketUrl: String? = null
)

