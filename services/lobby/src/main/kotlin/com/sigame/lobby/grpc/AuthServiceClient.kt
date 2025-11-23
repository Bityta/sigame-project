package com.sigame.lobby.grpc

import auth.AuthServiceGrpcKt
import auth.Auth
import com.sigame.lobby.config.AuthServiceConfig
import com.sigame.lobby.metrics.LobbyMetrics
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import io.grpc.Status
import io.grpc.StatusException
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.withContext
import mu.KotlinLogging
import org.springframework.stereotype.Service
import java.util.UUID
import java.util.concurrent.TimeUnit

private val logger = KotlinLogging.logger {}

data class UserInfo(
    val userId: UUID,
    val username: String
)

/**
 * Клиент для взаимодействия с Auth Service через gRPC
 * Поддерживает retry с exponential backoff
 */
@Service
class AuthServiceClient(
    private val config: AuthServiceConfig,
    private val lobbyMetrics: LobbyMetrics
) {
    
    private lateinit var channel: ManagedChannel
    private lateinit var stub: AuthServiceGrpcKt.AuthServiceCoroutineStub
    
    companion object {
        private const val MAX_RETRIES = 3
        private const val INITIAL_BACKOFF_MS = 100L
        private const val MAX_BACKOFF_MS = 2000L
        private const val TIMEOUT_SECONDS = 5L
    }
    
    @PostConstruct
    fun init() {
        channel = ManagedChannelBuilder
            .forAddress(config.host, config.port)
            .usePlaintext()
            .keepAliveTime(30, TimeUnit.SECONDS)
            .keepAliveTimeout(10, TimeUnit.SECONDS)
            .build()
        
        stub = AuthServiceGrpcKt.AuthServiceCoroutineStub(channel)
            .withDeadlineAfter(TIMEOUT_SECONDS, TimeUnit.SECONDS)
        
        logger.info { "Auth Service gRPC client initialized: ${config.host}:${config.port}" }
    }
    
    @PreDestroy
    fun shutdown() {
        if (::channel.isInitialized) {
            channel.shutdown()
            try {
                if (!channel.awaitTermination(5, TimeUnit.SECONDS)) {
                    channel.shutdownNow()
                }
            } catch (e: InterruptedException) {
                channel.shutdownNow()
                Thread.currentThread().interrupt()
            }
        }
    }
    
    /**
     * Валидирует JWT токен
     */
    suspend fun validateToken(token: String): UserInfo? = withRetry("validateToken") {
        withContext(Dispatchers.IO) {
            logger.debug { "Validating token with Auth Service..." }
            
            val request = Auth.ValidateTokenRequest.newBuilder()
                .setToken(token)
                .build()
            
            val response = stub.validateToken(request)
            
            if (response.valid) {
                UserInfo(
                    userId = UUID.fromString(response.userId),
                    username = response.username
                )
            } else {
                logger.warn { "Token validation failed: ${response.error}" }
                null
            }
        }
    }
    
    /**
     * Получает информацию о пользователе
     */
    suspend fun getUserInfo(userId: UUID): UserInfo? = withRetry("getUserInfo") {
        withContext(Dispatchers.IO) {
            val request = Auth.GetUserInfoRequest.newBuilder()
                .setUserId(userId.toString())
                .build()
            
            val response = stub.getUserInfo(request)
            
            if (response.error.isEmpty()) {
                UserInfo(
                    userId = UUID.fromString(response.userId),
                    username = response.username
                )
            } else {
                logger.warn { "Get user info failed: ${response.error}" }
                null
            }
        }
    }
    
    /**
     * Выполняет операцию с retry и exponential backoff
     */
    private suspend fun <T> withRetry(
        operationName: String,
        block: suspend () -> T
    ): T? {
        var lastException: Exception? = null
        var backoffMs = INITIAL_BACKOFF_MS
        
        repeat(MAX_RETRIES) { attempt ->
            try {
                val result = block()
                lobbyMetrics.recordGrpcCall("auth-service")
                return result
            } catch (e: StatusException) {
                lastException = e
                
                val shouldRetry = when (e.status.code) {
                    Status.Code.UNAVAILABLE,
                    Status.Code.DEADLINE_EXCEEDED,
                    Status.Code.RESOURCE_EXHAUSTED -> true
                    else -> false
                }
                
                if (shouldRetry && attempt < MAX_RETRIES - 1) {
                    logger.warn { "gRPC call to Auth Service failed (attempt ${attempt + 1}/$MAX_RETRIES): ${e.status}. Retrying in ${backoffMs}ms..." }
                    delay(backoffMs)
                    backoffMs = (backoffMs * 2).coerceAtMost(MAX_BACKOFF_MS)
                } else {
                    logger.error(e) { "gRPC call to Auth Service failed: $operationName" }
                    lobbyMetrics.recordGrpcError("auth-service", e.status.code.name)
                    return null
                }
            } catch (e: Exception) {
                logger.error(e) { "Unexpected error in gRPC call to Auth Service: $operationName" }
                lobbyMetrics.recordGrpcError("auth-service", "UNKNOWN")
                lastException = e
                return null
            }
        }
        
        return null
    }
}

