package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonProperty
import java.util.UUID

data class PlayerDto(
    @JsonProperty("userId")
    val userId: UUID,

    @JsonProperty("username")
    val username: String,

    @JsonProperty("avatar_url")
    val avatarUrl: String?,

    @JsonProperty("role")
    val role: String
)
