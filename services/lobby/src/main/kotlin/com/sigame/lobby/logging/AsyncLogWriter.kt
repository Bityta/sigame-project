package com.sigame.lobby.logging

import com.fasterxml.jackson.databind.ObjectMapper
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.cancel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.launch
import mu.KotlinLogging
import org.springframework.http.server.reactive.ServerHttpRequest
import org.springframework.http.server.reactive.ServerHttpResponse
import org.springframework.stereotype.Component

private val logger = KotlinLogging.logger {}

enum class LogLevel { INFO, WARN, ERROR }

data class LogEntry(
    val data: Map<String, Any?>,
    val level: LogLevel
)

@Component
class AsyncLogWriter(
    private val objectMapper: ObjectMapper,
    private val sanitizer: LogSanitizer
) {
    private val scope = CoroutineScope(Dispatchers.IO + SupervisorJob())
    private val logChannel = Channel<LogEntry>(capacity = 1000)

    @PostConstruct
    fun start() {
        scope.launch {
            for (entry in logChannel) {
                writeLog(entry)
            }
        }
        logger.info { "AsyncLogWriter started" }
    }

    @PreDestroy
    fun stop() {
        logChannel.close()
        scope.cancel()
        logger.info { "AsyncLogWriter stopped" }
    }

    fun logRequest(request: ServerHttpRequest, requestId: String) {
        val entry = buildMap {
            put("event", "http_request_received")
            put("request_id", requestId)
            put("method", request.method.name())
            put("path", request.path.value())
            put("query_params", sanitizer.sanitizeQueryParams(request.queryParams.toSingleValueMap()))
            put("headers", sanitizer.sanitizeHeaders(request.headers))
            put("remote_address", request.remoteAddress?.address?.hostAddress ?: "unknown")
        }
        send(LogEntry(entry, LogLevel.INFO))
    }

    fun logResponse(
        request: ServerHttpRequest,
        response: ServerHttpResponse,
        requestId: String,
        requestBody: String,
        responseBody: String,
        durationMs: Long,
        error: Throwable?
    ) {
        val statusCode = response.statusCode?.value() ?: 0

        val entry = buildMap {
            put("event", "http_request_completed")
            put("request_id", requestId)
            put("method", request.method.name())
            put("path", request.path.value())
            put("status", statusCode)
            put("duration_ms", durationMs)

            if (requestBody.isNotBlank()) {
                put("request_body", sanitizer.truncateAndSanitize(requestBody))
            }

            if (responseBody.isNotBlank() && (statusCode >= 400 || logger.isDebugEnabled)) {
                put("response_body", sanitizer.truncateAndSanitize(responseBody))
            }

            error?.let {
                put("error", it.javaClass.simpleName)
                put("error_message", it.message ?: "Unknown error")
            }
        }

        val level = when {
            error != null -> LogLevel.ERROR
            statusCode >= 500 -> LogLevel.ERROR
            statusCode >= 400 -> LogLevel.WARN
            else -> LogLevel.INFO
        }

        send(LogEntry(entry, level))
    }

    private fun send(entry: LogEntry) {
        logChannel.trySend(entry)
    }

    private fun writeLog(entry: LogEntry) {
        val json = toJson(entry.data)
        when (entry.level) {
            LogLevel.INFO -> logger.info { json }
            LogLevel.WARN -> logger.warn { json }
            LogLevel.ERROR -> logger.error { json }
        }
    }

    private fun toJson(data: Map<String, Any?>): String =
        runCatching { objectMapper.writeValueAsString(data) }.getOrElse { data.toString() }
}

