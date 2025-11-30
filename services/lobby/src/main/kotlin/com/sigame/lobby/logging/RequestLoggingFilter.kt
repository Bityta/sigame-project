package com.sigame.lobby.logging

import com.fasterxml.jackson.databind.ObjectMapper
import com.sigame.lobby.controller.ApiRoutes
import mu.KotlinLogging
import org.slf4j.MDC
import org.springframework.core.Ordered
import org.springframework.core.annotation.Order
import org.springframework.core.io.buffer.DataBuffer
import org.springframework.http.HttpHeaders
import org.springframework.http.server.reactive.ServerHttpRequestDecorator
import org.springframework.http.server.reactive.ServerHttpResponseDecorator
import org.springframework.stereotype.Component
import org.springframework.web.server.ServerWebExchange
import org.springframework.web.server.WebFilter
import org.springframework.web.server.WebFilterChain
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import java.nio.charset.StandardCharsets
import java.util.UUID
import java.util.concurrent.atomic.AtomicReference

private val logger = KotlinLogging.logger {}

@Component
@Order(Ordered.HIGHEST_PRECEDENCE + 1)
class RequestLoggingFilter(
    private val objectMapper: ObjectMapper
) : WebFilter {

    companion object {
        private val SKIP_PATHS = ApiRoutes.SkipLogging.PATHS
        private val SENSITIVE_FIELDS = setOf("password", "token", "accessToken", "refreshToken", "authorization")
        private val SENSITIVE_HEADERS = setOf("authorization", "x-api-key", "cookie")
        private const val MAX_BODY_LOG_SIZE = 10_000 // Ограничение на размер тела для логирования
    }

    override fun filter(exchange: ServerWebExchange, chain: WebFilterChain): Mono<Void> {
        val path = exchange.request.path.value()
        
        // Пропускаем служебные эндпоинты
        if (SKIP_PATHS.any { path.contains(it, ignoreCase = true) }) {
            return chain.filter(exchange)
        }

        val requestId = exchange.request.headers.getFirst("X-Request-ID") 
            ?: UUID.randomUUID().toString()
        val startTime = System.currentTimeMillis()
        
        // Устанавливаем MDC контекст
        setupMdc(exchange, requestId)
        
        val requestBody = AtomicReference<String>("")
        val responseBody = AtomicReference<String>("")
        
        // Декорируем запрос для перехвата тела
        val decoratedExchange = decorateExchange(exchange, requestBody, responseBody)
        
        // Логируем входящий запрос
        logIncomingRequest(exchange, requestId)
        
        return chain.filter(decoratedExchange)
            .doOnSuccess {
                logCompletedRequest(
                    exchange = decoratedExchange,
                    requestId = requestId,
                    requestBody = requestBody.get(),
                    responseBody = responseBody.get(),
                    duration = System.currentTimeMillis() - startTime,
                    error = null
                )
            }
            .doOnError { error ->
                logCompletedRequest(
                    exchange = decoratedExchange,
                    requestId = requestId,
                    requestBody = requestBody.get(),
                    responseBody = responseBody.get(),
                    duration = System.currentTimeMillis() - startTime,
                    error = error
                )
            }
            .doFinally {
                clearMdc()
            }
    }

    private fun setupMdc(exchange: ServerWebExchange, requestId: String) {
        MDC.put("request_id", requestId)
        MDC.put("trace_id", requestId)
        MDC.put("http_method", exchange.request.method.name())
        MDC.put("http_path", exchange.request.path.value())
        
        exchange.request.headers.getFirst("X-User-ID")?.let {
            MDC.put("user_id", it)
        }
    }

    private fun clearMdc() {
        MDC.remove("request_id")
        MDC.remove("trace_id")
        MDC.remove("http_method")
        MDC.remove("http_path")
        MDC.remove("user_id")
    }

    private fun decorateExchange(
        exchange: ServerWebExchange,
        requestBody: AtomicReference<String>,
        responseBody: AtomicReference<String>
    ): ServerWebExchange {
        val decoratedRequest = object : ServerHttpRequestDecorator(exchange.request) {
            override fun getBody(): Flux<DataBuffer> {
                return super.getBody().doOnNext { buffer ->
                    captureBody(buffer, requestBody)
                }
            }
        }

        val decoratedResponse = object : ServerHttpResponseDecorator(exchange.response) {
            override fun writeWith(body: org.reactivestreams.Publisher<out DataBuffer>): Mono<Void> {
                return super.writeWith(
                    Flux.from(body).doOnNext { buffer ->
                        captureBody(buffer, responseBody)
                    }
                )
            }
        }

        return exchange.mutate()
            .request(decoratedRequest)
            .response(decoratedResponse)
            .build()
    }

    private fun captureBody(buffer: DataBuffer, target: AtomicReference<String>) {
        val readableBytes = buffer.readableByteCount()
        if (readableBytes > 0 && target.get().length < MAX_BODY_LOG_SIZE) {
            val bytes = ByteArray(minOf(readableBytes, MAX_BODY_LOG_SIZE - target.get().length))
            buffer.toByteBuffer().get(bytes)
            target.updateAndGet { it + String(bytes, StandardCharsets.UTF_8) }
        }
    }

    private fun logIncomingRequest(exchange: ServerWebExchange, requestId: String) {
        val request = exchange.request
        
        val logEntry = buildMap {
            put("event", "http_request_received")
            put("request_id", requestId)
            put("method", request.method.name())
            put("path", request.path.value())
            put("query_params", sanitizeQueryParams(request.queryParams.toSingleValueMap()))
            put("headers", sanitizeHeaders(request.headers))
            put("remote_address", request.remoteAddress?.address?.hostAddress ?: "unknown")
        }
        
        logger.info { toJson(logEntry) }
    }

    private fun logCompletedRequest(
        exchange: ServerWebExchange,
        requestId: String,
        requestBody: String,
        responseBody: String,
        duration: Long,
        error: Throwable?
    ) {
        val request = exchange.request
        val response = exchange.response
        val statusCode = response.statusCode?.value() ?: 0
        
        val logEntry = buildMap {
            put("event", "http_request_completed")
            put("request_id", requestId)
            put("method", request.method.name())
            put("path", request.path.value())
            put("status", statusCode)
            put("duration_ms", duration)
            
            // Тело запроса (если есть и не пустое)
            if (requestBody.isNotBlank()) {
                put("request_body", truncateAndSanitize(requestBody))
            }
            
            // Тело ответа для ошибок или debug уровня
            if (responseBody.isNotBlank() && (statusCode >= 400 || logger.isDebugEnabled)) {
                put("response_body", truncateAndSanitize(responseBody))
            }
            
            // Информация об ошибке
            error?.let {
                put("error", it.javaClass.simpleName)
                put("error_message", it.message ?: "Unknown error")
            }
        }
        
        when {
            error != null -> logger.error { toJson(logEntry) }
            statusCode >= 500 -> logger.error { toJson(logEntry) }
            statusCode >= 400 -> logger.warn { toJson(logEntry) }
            else -> logger.info { toJson(logEntry) }
        }
    }

    private fun sanitizeHeaders(headers: HttpHeaders): Map<String, String> {
        return headers.toSingleValueMap()
            .filterKeys { key -> 
                !SENSITIVE_HEADERS.any { it.equals(key, ignoreCase = true) }
            }
            .mapValues { (key, value) ->
                if (SENSITIVE_HEADERS.any { it.equals(key, ignoreCase = true) }) {
                    "***HIDDEN***"
                } else {
                    value
                }
            }
    }

    private fun sanitizeQueryParams(params: Map<String, String>): Map<String, String> {
        return params.mapValues { (key, value) ->
            if (SENSITIVE_FIELDS.any { it.equals(key, ignoreCase = true) }) {
                "***HIDDEN***"
            } else {
                value
            }
        }
    }

    private fun truncateAndSanitize(body: String): String {
        val truncated = if (body.length > MAX_BODY_LOG_SIZE) {
            body.take(MAX_BODY_LOG_SIZE) + "...[TRUNCATED]"
        } else {
            body
        }
        return sanitizeBody(truncated)
    }

    private fun sanitizeBody(body: String): String {
        if (body.isBlank()) return body
        
        var sanitized = body
        SENSITIVE_FIELDS.forEach { field ->
            // Маскируем JSON поля
            sanitized = sanitized.replace(
                Regex(""""$field"\s*:\s*"[^"]*"""", RegexOption.IGNORE_CASE),
                """"$field":"***HIDDEN***""""
            )
            // Маскируем query параметры
            sanitized = sanitized.replace(
                Regex("""$field=[^&\s]*""", RegexOption.IGNORE_CASE),
                "$field=***HIDDEN***"
            )
        }
        return sanitized
    }

    private fun toJson(data: Map<String, Any?>): String {
        return try {
            objectMapper.writeValueAsString(data)
        } catch (e: Exception) {
            data.toString()
        }
    }
}
