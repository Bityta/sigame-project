package com.sigame.lobby.service

import com.sigame.lobby.config.RoomConfig
import com.sigame.lobby.domain.GameRoomRepository
import kotlinx.coroutines.reactive.awaitFirstOrNull
import org.springframework.stereotype.Service
import kotlin.random.Random

@Service
class RoomCodeGenerator(
    private val roomConfig: RoomConfig,
    private val gameRoomRepository: GameRoomRepository
) {
    
    private val charset = roomConfig.codeCharset.toCharArray()
    
    suspend fun generateUniqueCode(): String {
        var attempts = 0
        val maxAttempts = 100
        
        while (attempts < maxAttempts) {
            val code = generateCode()
            val existing = gameRoomRepository.findByRoomCode(code).awaitFirstOrNull()
            
            if (existing == null) {
                return code
            }
            
            attempts++
        }
        
        throw IllegalStateException("Failed to generate unique room code after $maxAttempts attempts")
    }
    
    private fun generateCode(): String {
        return (1..roomConfig.codeLength)
            .map { charset[Random.nextInt(charset.size)] }
            .joinToString("")
    }
}

