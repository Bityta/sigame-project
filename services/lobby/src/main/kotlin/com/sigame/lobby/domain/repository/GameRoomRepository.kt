package com.sigame.lobby.domain.repository

import com.sigame.lobby.domain.model.GameRoom
import org.springframework.data.r2dbc.repository.Query
import org.springframework.data.r2dbc.repository.R2dbcRepository
import org.springframework.stereotype.Repository
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import java.util.UUID

@Repository
interface GameRoomRepository : R2dbcRepository<GameRoom, UUID> {
    
    fun findByRoomCode(roomCode: String): Mono<GameRoom>
    
    @Query("""
        SELECT * FROM game_rooms 
        WHERE status = :status 
        ORDER BY created_at DESC 
        LIMIT :limit OFFSET :offset
    """)
    fun findByStatus(status: String, limit: Int, offset: Int): Flux<GameRoom>
    
    @Query("""
        SELECT * FROM game_rooms 
        WHERE is_public = true AND status = 'waiting'
        ORDER BY created_at DESC 
        LIMIT :limit OFFSET :offset
    """)
    fun findPublicWaitingRooms(limit: Int, offset: Int): Flux<GameRoom>
    
    @Query("SELECT COUNT(*) FROM game_rooms WHERE status = :status")
    fun countByStatus(status: String): Mono<Long>
    
    @Query("SELECT COUNT(*) FROM game_rooms WHERE is_public = true AND status = 'waiting'")
    fun countPublicWaitingRooms(): Mono<Long>
    
    fun findByHostId(hostId: UUID): Flux<GameRoom>
}

