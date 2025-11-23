package com.sigame.lobby.service.external

import com.sigame.lobby.config.GameServiceConfig
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.grpc.AuthServiceClient
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import mu.KotlinLogging
import org.springframework.http.MediaType
import org.springframework.stereotype.Service
import org.springframework.web.reactive.function.client.WebClient
import org.springframework.web.reactive.function.client.awaitBody
import java.util.UUID

private val logger = KotlinLogging.logger {}

/**
 * Клиент для взаимодействия с Game Service
 */
@Service
class GameServiceClient(
    private val gameServiceConfig: GameServiceConfig,
    private val authServiceClient: AuthServiceClient,
    private val webClient: WebClient = WebClient.builder().build()
) {
    
    /**
     * Создает игровую сессию в Game Service
     */
    suspend fun createGameSession(
        roomId: UUID,
        packId: UUID,
        players: List<RoomPlayer>,
        settings: GameSettings
    ): GameSessionResponse = withContext(Dispatchers.IO) {
        try {
            val createGameUrl = "${gameServiceConfig.baseUrl}/api/game/create"
            logger.info { "Creating game session at $createGameUrl for room $roomId" }
            
            // Подготовка данных игроков с их username
            val playerData = players.map { player ->
                val userInfo = authServiceClient.getUserInfo(player.userId)
                mapOf(
                    "user_id" to player.userId.toString(),
                    "username" to (userInfo?.username ?: "Unknown"),
                    "role" to player.role
                )
            }
            
            val request = mapOf(
                "room_id" to roomId.toString(),
                "pack_id" to packId.toString(),
                "players" to playerData,
                "settings" to mapOf(
                    "time_for_answer" to settings.timeForAnswer,
                    "time_for_choice" to settings.timeForChoice,
                    "allow_wrong_answer" to settings.allowWrongAnswer,
                    "show_right_answer" to settings.showRightAnswer
                )
            )
            
            val response = webClient.post()
                .uri(createGameUrl)
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(request)
                .retrieve()
                .awaitBody<Map<String, Any>>()
            
            val gameSessionId = response["game_id"] as? String 
                ?: throw IllegalStateException("Game Service did not return game_id")
            val wsUrl = response["websocket_url"] as? String 
                ?: throw IllegalStateException("Game Service did not return websocket_url")
            
            logger.info { "Game session created: $gameSessionId with WebSocket at $wsUrl" }
            
            GameSessionResponse(gameSessionId, wsUrl)
            
        } catch (e: Exception) {
            logger.error(e) { "Failed to create game session: ${e.message}" }
            throw IllegalStateException("Failed to start game session: ${e.message}", e)
        }
    }
}

/**
 * Настройки игры
 */
data class GameSettings(
    val timeForAnswer: Int = 30,
    val timeForChoice: Int = 60,
    val allowWrongAnswer: Boolean = true,
    val showRightAnswer: Boolean = true
)

/**
 * Ответ от Game Service при создании сессии
 */
data class GameSessionResponse(
    val gameSessionId: String,
    val wsUrl: String
)

