package com.sigame.lobby.domain.model

import com.sigame.lobby.domain.enums.RoomStatus
import org.springframework.data.annotation.Id
import org.springframework.data.annotation.Transient
import org.springframework.data.domain.Persistable
import org.springframework.data.relational.core.mapping.Column
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime
import java.util.UUID

@Table("game_rooms")
data class GameRoom(
    @Id
    @get:JvmName("getId_")
    val id: UUID = UUID.randomUUID(),
    @Column("room_code")
    val roomCode: String,
    @Column("host_id")
    val hostId: UUID,
    @Column("pack_id")
    val packId: UUID,
    @Column("name")
    val name: String,
    @Column("status")
    val status: String = "waiting",
    @Column("max_players")
    val maxPlayers: Int = 6,
    @Column("is_public")
    val isPublic: Boolean = true,
    @Column("password_hash")
    val passwordHash: String? = null,
    @Column("created_at")
    val createdAt: LocalDateTime = LocalDateTime.now(),
    @Column("updated_at")
    val updatedAt: LocalDateTime = LocalDateTime.now(),
    @Column("started_at")
    val startedAt: LocalDateTime? = null,
    @Column("finished_at")
    val finishedAt: LocalDateTime? = null,
    @Transient
    private val isNewEntity: Boolean = true
) : Persistable<UUID> {
    
    override fun getId(): UUID = id
    override fun isNew(): Boolean = isNewEntity
    
    fun getStatusEnum(): RoomStatus = RoomStatus.valueOf(status.uppercase())
    
    fun markPersisted(): GameRoom = copy(isNewEntity = false)

    companion object {
        fun statusFromEnum(status: RoomStatus): String = status.name.lowercase()
    }
}

