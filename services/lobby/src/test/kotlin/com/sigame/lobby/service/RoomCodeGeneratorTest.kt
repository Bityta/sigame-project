package com.sigame.lobby.service

import com.sigame.lobby.config.RoomConfig
import com.sigame.lobby.domain.repository.GameRoomRepository
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.mockito.kotlin.mock
import org.mockito.kotlin.whenever
import reactor.core.publisher.Mono

class RoomCodeGeneratorTest {

    private lateinit var roomConfig: RoomConfig
    private lateinit var gameRoomRepository: GameRoomRepository
    private lateinit var roomCodeGenerator: RoomCodeGenerator

    @BeforeEach
    fun setup() {
        roomConfig = RoomConfig(
            codeLength = 6,
            codeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
        )
        gameRoomRepository = mock()
        roomCodeGenerator = RoomCodeGenerator(roomConfig, gameRoomRepository)
    }

    @Test
    fun `should generate code with correct length`() = runBlocking {
        whenever(gameRoomRepository.existsByCode(org.mockito.kotlin.any())).thenReturn(Mono.just(false))

        val code = roomCodeGenerator.generateUniqueCode()
        
        assertEquals(6, code.length)
    }

    @Test
    fun `should generate code with valid characters`() = runBlocking {
        whenever(gameRoomRepository.existsByCode(org.mockito.kotlin.any())).thenReturn(Mono.just(false))

        val code = roomCodeGenerator.generateUniqueCode()
        
        assertTrue(code.all { it in roomConfig.codeCharset })
    }

    @Test
    fun `should generate unique codes`() = runBlocking {
        whenever(gameRoomRepository.existsByCode(org.mockito.kotlin.any())).thenReturn(Mono.just(false))

        val codes = mutableSetOf<String>()
        repeat(100) {
            codes.add(roomCodeGenerator.generateUniqueCode())
        }
        
        // All codes should be unique
        assertEquals(100, codes.size)
    }
}

