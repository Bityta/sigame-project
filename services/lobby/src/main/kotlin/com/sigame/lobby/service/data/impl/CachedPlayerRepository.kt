package com.sigame.lobby.service.data.impl

import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.service.cache.RoomCacheService
import com.sigame.lobby.service.data.PlayerRepository
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.launch
import kotlinx.coroutines.reactive.asFlow
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import org.springframework.stereotype.Repository
import java.util.UUID

@Repository
class CachedPlayerRepository(
    private val dbRepository: RoomPlayerRepository,
    private val cache: RoomCacheService
) : PlayerRepository {

    override suspend fun findActiveByUserId(userId: UUID): RoomPlayer? {
        cache.getUserCurrentRoom(userId)?.let { cachedRoomId ->
            val player = dbRepository.findByRoomIdAndUserId(cachedRoomId, userId).awaitFirstOrNull()
            if (player != null && player.leftAt == null) return player
            CoroutineScope(Dispatchers.IO).launch { cache.deleteUserCurrentRoom(userId) }
        }

        return dbRepository.findActiveByUserId(userId).awaitFirstOrNull()?.also { player ->
            CoroutineScope(Dispatchers.IO).launch {
                cache.setUserCurrentRoom(userId, player.roomId)
            }
        }
    }

    override suspend fun findByRoomAndUser(roomId: UUID, userId: UUID): RoomPlayer? {
        return dbRepository.findByRoomIdAndUserId(roomId, userId).awaitFirstOrNull()
    }

    override suspend fun findActiveByRoomId(roomId: UUID): List<RoomPlayer> {
        return dbRepository.findActiveByRoomId(roomId).asFlow().toList()
    }

    override suspend fun findActiveByRoomIds(roomIds: List<UUID>): Map<UUID, List<RoomPlayer>> {
        if (roomIds.isEmpty()) return emptyMap()
        return dbRepository.findActiveByRoomIds(roomIds).asFlow().toList().groupBy { it.roomId }
    }

    override suspend fun countActiveByRoomId(roomId: UUID): Int {
        return dbRepository.countActiveByRoomId(roomId).awaitFirstOrNull()?.toInt() ?: 0
    }

    override suspend fun save(player: RoomPlayer): RoomPlayer {
        val saved = dbRepository.save(player).awaitFirst()

        CoroutineScope(Dispatchers.IO).launch {
            if (saved.leftAt == null) {
                cache.setUserCurrentRoom(saved.userId, saved.roomId)
            } else {
                cache.deleteUserCurrentRoom(saved.userId)
            }
        }

        return saved
    }
}

