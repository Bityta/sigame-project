package com.sigame.lobby.domain.model

import org.springframework.data.annotation.Id
import org.springframework.data.annotation.Transient
import org.springframework.data.domain.Persistable
import org.springframework.data.relational.core.mapping.Column
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime
import java.util.UUID

@Table("room_settings")
class RoomSettings(
    @Id
    @Column("room_id")
    val roomId: UUID,
    @Column("time_for_answer")
    val timeForAnswer: Int = 30,
    @Column("time_for_choice")
    val timeForChoice: Int = 60,
    @Column("allow_wrong_answer")
    val allowWrongAnswer: Boolean = true,
    @Column("show_right_answer")
    val showRightAnswer: Boolean = true,
    @Column("created_at")
    val createdAt: LocalDateTime = LocalDateTime.now(),
    @Column("updated_at")
    val updatedAt: LocalDateTime = LocalDateTime.now()
) : Persistable<UUID> {
    
    @Transient
    private var isNewEntity: Boolean = true
    
    override fun getId(): UUID = roomId
    override fun isNew(): Boolean = isNewEntity
    
    fun markPersisted(): RoomSettings {
        isNewEntity = false
        return this
    }
    
    fun copy(
        roomId: UUID = this.roomId,
        timeForAnswer: Int = this.timeForAnswer,
        timeForChoice: Int = this.timeForChoice,
        allowWrongAnswer: Boolean = this.allowWrongAnswer,
        showRightAnswer: Boolean = this.showRightAnswer,
        createdAt: LocalDateTime = this.createdAt,
        updatedAt: LocalDateTime = this.updatedAt
    ): RoomSettings {
        val copy = RoomSettings(roomId, timeForAnswer, timeForChoice, allowWrongAnswer, showRightAnswer, createdAt, updatedAt)
        copy.isNewEntity = this.isNewEntity
        return copy
    }
}
