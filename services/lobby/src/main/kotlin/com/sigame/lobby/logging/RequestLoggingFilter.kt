package com.sigame.lobby.logging

import kotlinx.coroutines.reactor.mono
import mu.KotlinLogging
import org.slf4j.MDC
import org.springframework.core.io.buffer.DataBuffer
import org.springframework.core.io.buffer.DataBufferUtils
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

private val logger = KotlinLogging.logger {}

/**
 * Фильтр для логирования запросов и ответов с телами (асинхронно)
 */
@Component
class RequestLoggingFilter : WebFilter {
    
    private val sensitiveFields = setOf("password", "token", "accessToken", "refreshToken")
    
    override fun filter(exchange: ServerWebExchange, chain: WebFilterChain): Mono<Void> {
        val requestId = exchange.request.headers.getFirst("X-Request-ID")
            ?: UUID.randomUUID().toString()
        
        val userId = exchange.request.headers.getFirst("X-User-ID")
        val path = exchange.request.path.value()
        val method = exchange.request.method.name()
        
        // Skip health and metrics endpoints
        if (path.contains("/health") || path.contains("/metrics") || path.contains("/actuator")) {
            return chain.filter(exchange)
        }
        
        // Добавляем в MDC
        MDC.put("request_id", requestId)
        MDC.put("trace_id", requestId)
        if (userId != null) {
            MDC.put("user_id", userId)
        }
        MDC.put("request_path", path)
        MDC.put("request_method", method)
        
        val startTime = System.currentTimeMillis()
        
        // Wrap request to log body
        val cachedBodyExchange = if (exchange.request.headers.contentLength > 0) {
            val cachedBody = StringBuilder()
            val decoratedRequest = object : ServerHttpRequestDecorator(exchange.request) {
                override fun getBody(): Flux<DataBuffer> {
                    return super.getBody().doOnNext { buffer ->
                        // Read body asynchronously
                        mono {
                            val bytes = ByteArray(buffer.readableByteCount())
                            buffer.read(bytes)
                            DataBufferUtils.release(buffer)
                            val bodyString = String(bytes, StandardCharsets.UTF_8)
                            cachedBody.append(bodyString)
                        }.subscribe()
                    }.doOnComplete {
                        // Log request with body
                        logger.debug { "Request body: ${sanitizeBody(cachedBody.toString())}" }
                    }
                }
            }
            exchange.mutate().request(decoratedRequest).build()
        } else {
            logger.debug { "Incoming request: $method $path (no body)" }
            exchange
        }
        
        // Wrap response to log body
        val responseBody = StringBuilder()
        val decoratedResponse = object : ServerHttpResponseDecorator(cachedBodyExchange.response) {
            override fun writeWith(body: org.reactivestreams.Publisher<out DataBuffer>): Mono<Void> {
                return super.writeWith(
                    Flux.from(body).doOnNext { buffer ->
                        // Read response body asynchronously
                        mono {
                            val bytes = ByteArray(buffer.readableByteCount())
                            buffer.slice().read(bytes)
                            val bodyString = String(bytes, StandardCharsets.UTF_8)
                            responseBody.append(bodyString)
                        }.subscribe()
                    }
                )
            }
        }
        
        val responseExchange = cachedBodyExchange.mutate().response(decoratedResponse).build()
        
        return chain.filter(responseExchange)
            .doFinally {
                val duration = System.currentTimeMillis() - startTime
                val statusCode = responseExchange.response.statusCode?.value() ?: 0
                
                // Log response asynchronously
                mono {
                    if (responseBody.isNotEmpty()) {
                        logger.debug { "Response body: ${sanitizeBody(responseBody.toString())}" }
                    }
                    logger.info { "Request completed: $method $path - status=$statusCode duration=${duration}ms" }
                }.subscribe()
                
                // Очищаем MDC
                MDC.remove("request_id")
                MDC.remove("trace_id")
                MDC.remove("user_id")
                MDC.remove("request_path")
                MDC.remove("request_method")
            }
    }
    
    private fun sanitizeBody(body: String): String {
        if (body.isBlank()) return body
        
        var sanitized = body
        sensitiveFields.forEach { field ->
            // Simple regex to hide sensitive values
            sanitized = sanitized.replace(
                Regex(""""$field"\s*:\s*"[^"]*""""),
                """"$field":"***HIDDEN***""""
            )
        }
        return sanitized
    }
}

