package com.sigame.lobby.domain.repository

import com.sigame.lobby.domain.model.RoomPlayer
import org.springframework.data.r2dbc.repository.Query
import org.springframework.data.r2dbc.repository.R2dbcRepository
import org.springframework.stereotype.Repository
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import java.util.UUID

/**
 * Репозиторий для работы с игроками в комнатах
 */
@Repository
interface RoomPlayerRepository : R2dbcRepository<RoomPlayer, UUID> {
    
    fun findByRoomId(roomId: UUID): Flux<RoomPlayer>
    
    @Query("SELECT * FROM room_players WHERE room_id = :roomId AND left_at IS NULL")
    fun findActiveByRoomId(roomId: UUID): Flux<RoomPlayer>
    
    @Query("SELECT COUNT(*) FROM room_players WHERE room_id = :roomId AND left_at IS NULL")
    fun countActiveByRoomId(roomId: UUID): Mono<Long>
    
    fun findByUserId(userId: UUID): Flux<RoomPlayer>
    
    @Query("SELECT * FROM room_players WHERE user_id = :userId AND left_at IS NULL")
    fun findActiveByUserId(userId: UUID): Mono<RoomPlayer>
    
    @Query("SELECT * FROM room_players WHERE room_id = :roomId AND user_id = :userId")
    fun findByRoomIdAndUserId(roomId: UUID, userId: UUID): Mono<RoomPlayer>
    
    @Query("SELECT * FROM room_players WHERE left_at IS NULL")
    fun findByLeftAtIsNull(): Flux<RoomPlayer>
}

