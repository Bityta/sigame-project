package com.sigame.lobby.domain.model

import org.springframework.data.annotation.Id
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime
import java.util.UUID

@Table("room_settings")
data class RoomSettings(
    @Id
    val roomId: UUID,
    val timeForAnswer: Int = 30,
    val timeForChoice: Int = 60,
    val allowWrongAnswer: Boolean = true,
    val showRightAnswer: Boolean = true,
    val createdAt: LocalDateTime = LocalDateTime.now(),
    val updatedAt: LocalDateTime = LocalDateTime.now()
)

