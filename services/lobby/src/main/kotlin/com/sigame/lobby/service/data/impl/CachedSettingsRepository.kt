package com.sigame.lobby.service.data.impl

import com.sigame.lobby.domain.model.RoomSettings
import com.sigame.lobby.domain.repository.RoomSettingsRepository
import com.sigame.lobby.service.cache.RoomCacheService
import com.sigame.lobby.service.data.SettingsRepository
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import org.springframework.stereotype.Repository
import java.util.UUID

@Repository
class CachedSettingsRepository(
    private val dbRepository: RoomSettingsRepository,
    private val cache: RoomCacheService
) : SettingsRepository {

    override suspend fun findByRoomId(roomId: UUID): RoomSettings? {
        cache.getRoomSettings(roomId)?.let { return it }

        val settings = dbRepository.findByRoomId(roomId).awaitFirstOrNull()
        if (settings != null) {
            CoroutineScope(Dispatchers.IO).launch {
                cache.setRoomSettings(roomId, settings)
            }
        }
        return settings
    }

    override suspend fun save(settings: RoomSettings): RoomSettings {
        val saved = dbRepository.save(settings).awaitFirst()
        CoroutineScope(Dispatchers.IO).launch {
            cache.setRoomSettings(saved.roomId, saved)
        }
        return saved
    }

    override suspend fun insert(settings: RoomSettings): RoomSettings {
        dbRepository.insertRoomSettings(
            roomId = settings.roomId,
            timeForAnswer = settings.timeForAnswer,
            timeForChoice = settings.timeForChoice
        ).awaitFirstOrNull()

        CoroutineScope(Dispatchers.IO).launch {
            cache.setRoomSettings(settings.roomId, settings)
        }
        return settings
    }
}

