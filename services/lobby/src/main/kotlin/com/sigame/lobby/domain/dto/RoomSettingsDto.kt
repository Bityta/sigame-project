package com.sigame.lobby.domain.dto

import jakarta.validation.constraints.Max
import jakarta.validation.constraints.Min

/**
 * DTO настроек комнаты
 */
data class RoomSettingsDto(
    @field:Min(value = 10, message = "Minimum 10 seconds for answer")
    @field:Max(value = 120, message = "Maximum 120 seconds for answer")
    val timeForAnswer: Int = 30,
    
    @field:Min(value = 10, message = "Minimum 10 seconds for choice")
    @field:Max(value = 180, message = "Maximum 180 seconds for choice")
    val timeForChoice: Int = 60,
    
    val allowWrongAnswer: Boolean = true,
    val showRightAnswer: Boolean = true
)

