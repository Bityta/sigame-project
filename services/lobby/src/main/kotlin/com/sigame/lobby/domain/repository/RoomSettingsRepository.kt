package com.sigame.lobby.domain.repository

import com.sigame.lobby.domain.model.RoomSettings
import org.springframework.data.r2dbc.repository.Modifying
import org.springframework.data.r2dbc.repository.Query
import org.springframework.data.r2dbc.repository.R2dbcRepository
import org.springframework.stereotype.Repository
import reactor.core.publisher.Mono
import java.util.UUID

@Repository
interface RoomSettingsRepository : R2dbcRepository<RoomSettings, UUID> {
    
    fun findByRoomId(roomId: UUID): Mono<RoomSettings>
    
    @Modifying
    @Query("""
        INSERT INTO room_settings (room_id, time_for_answer, time_for_choice)
        VALUES (:roomId, :timeForAnswer, :timeForChoice)
    """)
    fun insertRoomSettings(
        roomId: UUID,
        timeForAnswer: Int,
        timeForChoice: Int
    ): Mono<Int>
}
