package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonIgnoreProperties
import com.fasterxml.jackson.annotation.JsonInclude
import com.fasterxml.jackson.annotation.JsonProperty

/**
 * Запрос на присоединение к комнате
 * Все поля необязательны
 */
@JsonIgnoreProperties(ignoreUnknown = true)
@JsonInclude(JsonInclude.Include.NON_NULL)
data class JoinRoomRequest(
    @JsonProperty("password")
    val password: String? = null,
    
    @JsonProperty("role")
    val role: String = "player"
)

