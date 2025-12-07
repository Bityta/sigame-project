package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty

data class SetReadyRequest(
    @JsonProperty("isReady")
    val isReady: Boolean = true
)

