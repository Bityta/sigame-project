package com.sigame.lobby.grpc.pack

import com.sigame.lobby.config.PackServiceConfig
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
import com.sigame.pack.proto.GetPackInfoRequest
import com.sigame.pack.proto.PackServiceGrpcKt
import com.sigame.pack.proto.ValidatePackRequest
import java.util.UUID
import java.util.concurrent.TimeUnit

private val logger = KotlinLogging.logger {}

@Service
class PackServiceClient(
    private val config: PackServiceConfig,
    private val lobbyMetrics: LobbyMetrics
) {

    private lateinit var channel: ManagedChannel
    private lateinit var stub: PackServiceGrpcKt.PackServiceCoroutineStub

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

        stub = PackServiceGrpcKt.PackServiceCoroutineStub(channel)

        logger.info { "Pack Service gRPC client initialized: ${config.host}:${config.port}" }
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

    suspend fun validatePack(packId: UUID, userId: UUID? = null): PackValidationResult? = withRetry("validatePack") {
        withContext(Dispatchers.IO) {
            val requestBuilder = ValidatePackRequest.newBuilder()
                .setPackId(packId.toString())

            userId?.let { requestBuilder.setUserId(it.toString()) }

            val response = stub
                .withDeadlineAfter(TIMEOUT_SECONDS, TimeUnit.SECONDS)
                .validatePackExists(requestBuilder.build())

            PackValidationResult(
                exists = response.exists,
                isOwner = response.isOwner,
                status = response.status.ifEmpty { "approved" },
                error = response.error.ifEmpty { null }
            )
        }
    }

    suspend fun getPackInfo(packId: UUID): PackInfo? = withRetry("getPackInfo") {
        withContext(Dispatchers.IO) {
            val request = GetPackInfoRequest.newBuilder()
                .setPackId(packId.toString())
                .build()

            val response = stub
                .withDeadlineAfter(TIMEOUT_SECONDS, TimeUnit.SECONDS)
                .getPackInfo(request)
            PackInfo(
                id = UUID.fromString(response.packId),
                name = response.name,
                author = response.author,
                roundsCount = response.roundsCount,
                questionsCount = response.questionsCount
            )
        }
    }

    private suspend fun <T> withRetry(
        operationName: String,
        block: suspend () -> T
    ): T? {
        var backoffMs = INITIAL_BACKOFF_MS

        repeat(MAX_RETRIES) { attempt ->
            try {
                val result = block()
                lobbyMetrics.recordGrpcCall("pack-service")
                return result
            } catch (e: StatusException) {
                val shouldRetry = when (e.status.code) {
                    Status.Code.UNAVAILABLE,
                    Status.Code.DEADLINE_EXCEEDED,
                    Status.Code.RESOURCE_EXHAUSTED -> true
                    else -> false
                }

                if (shouldRetry && attempt < MAX_RETRIES - 1) {
                    logger.warn { "gRPC call to Pack Service failed (attempt ${attempt + 1}/$MAX_RETRIES): ${e.status}. Retrying in ${backoffMs}ms..." }
                    delay(backoffMs)
                    backoffMs = (backoffMs * 2).coerceAtMost(MAX_BACKOFF_MS)
                } else {
                    logger.error(e) { "gRPC call to Pack Service failed: $operationName" }
                    lobbyMetrics.recordGrpcError("pack-service", e.status.code.name)
                    return null
                }
            } catch (e: Exception) {
                logger.error(e) { "Unexpected error in gRPC call to Pack Service: $operationName" }
                lobbyMetrics.recordGrpcError("pack-service", "UNKNOWN")
                return null
            }
        }

        return null
    }
}

