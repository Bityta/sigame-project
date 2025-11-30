package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonIgnoreProperties
import com.fasterxml.jackson.annotation.JsonProperty
import jakarta.validation.constraints.Max
import jakarta.validation.constraints.Min

@JsonIgnoreProperties(ignoreUnknown = true)
data class RoomSettingsDto(
    @field:Min(value = 10, message = "Minimum 10 seconds for answer")
    @field:Max(value = 120, message = "Maximum 120 seconds for answer")
    @JsonProperty("timeForAnswer")
    val timeForAnswer: Int = 30,
    
    @field:Min(value = 10, message = "Minimum 10 seconds for choice")
    @field:Max(value = 180, message = "Maximum 180 seconds for choice")
    @JsonProperty("timeForChoice")
    val timeForChoice: Int = 60,
    
    @JsonProperty("allowWrongAnswer")
    val allowWrongAnswer: Boolean = true,
    
    @JsonProperty("showRightAnswer")
    val showRightAnswer: Boolean = true
)
