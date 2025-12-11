package com.sigame.lobby.domain.model

import org.springframework.data.annotation.Id
import org.springframework.data.relational.core.mapping.Column
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime
import java.util.UUID

@Table("room_settings")
data class RoomSettings(
    @Id
    @Column("room_id")
    val roomId: UUID,
    @Column("time_for_answer")
    val timeForAnswer: Int = 30,
    @Column("time_for_choice")
    val timeForChoice: Int = 60,
    @Column("created_at")
    val createdAt: LocalDateTime = LocalDateTime.now(),
    @Column("updated_at")
    val updatedAt: LocalDateTime = LocalDateTime.now()
)
