package com.sigame.lobby.domain.repository

import com.sigame.lobby.domain.model.RoomSettings
import org.springframework.data.r2dbc.repository.Modifying
import org.springframework.data.r2dbc.repository.Query
import org.springframework.data.r2dbc.repository.R2dbcRepository
import org.springframework.stereotype.Repository
import reactor.core.publisher.Mono
import java.util.UUID

/**
 * Репозиторий для работы с настройками комнат
 */
@Repository
interface RoomSettingsRepository : R2dbcRepository<RoomSettings, UUID> {
    
    fun findByRoomId(roomId: UUID): Mono<RoomSettings>
    
    @Modifying
    @Query("""
        INSERT INTO room_settings (room_id, time_for_answer, time_for_choice, allow_wrong_answer, show_right_answer)
        VALUES (:roomId, :timeForAnswer, :timeForChoice, :allowWrongAnswer, :showRightAnswer)
    """)
    fun insertRoomSettings(
        roomId: UUID,
        timeForAnswer: Int,
        timeForChoice: Int,
        allowWrongAnswer: Boolean,
        showRightAnswer: Boolean
    ): Mono<Int>
}

