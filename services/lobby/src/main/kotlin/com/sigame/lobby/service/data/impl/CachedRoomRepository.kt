package com.sigame.lobby.service.data.impl

import com.sigame.lobby.domain.exception.RoomNotFoundException
import com.sigame.lobby.domain.exception.RoomNotFoundByCodeException
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.service.data.RoomRepository
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import org.springframework.stereotype.Repository
import java.util.UUID

@Repository
class CachedRoomRepository(
    private val dbRepository: GameRoomRepository
) : RoomRepository {

    override suspend fun findById(id: UUID): GameRoom {
        return dbRepository.findById(id).awaitFirstOrNull()
            ?: throw RoomNotFoundException(id)
    }

    override suspend fun findByCode(code: String): GameRoom {
        return dbRepository.findByRoomCode(code).awaitFirstOrNull()
            ?: throw RoomNotFoundByCodeException(code)
    }

    override suspend fun save(room: GameRoom): GameRoom {
        return dbRepository.save(room).awaitFirst()
    }
}

