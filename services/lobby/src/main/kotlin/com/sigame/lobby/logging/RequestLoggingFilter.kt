package com.sigame.lobby.logging

import kotlinx.coroutines.reactor.awaitSingleOrNull
import mu.KotlinLogging
import org.slf4j.MDC
import org.springframework.stereotype.Component
import org.springframework.web.server.ServerWebExchange
import org.springframework.web.server.WebFilter
import org.springframework.web.server.WebFilterChain
import reactor.core.publisher.Mono
import java.util.UUID

private val logger = KotlinLogging.logger {}

/**
 * Фильтр для добавления request ID и user ID в MDC контекст для логирования
 */
@Component
class RequestLoggingFilter : WebFilter {
    
    override fun filter(exchange: ServerWebExchange, chain: WebFilterChain): Mono<Void> {
        val requestId = exchange.request.headers.getFirst("X-Request-ID")
            ?: UUID.randomUUID().toString()
        
        val userId = exchange.request.headers.getFirst("X-User-ID")
        val path = exchange.request.path.value()
        val method = exchange.request.method.name()
        
        // Добавляем в MDC
        MDC.put("request_id", requestId)
        MDC.put("trace_id", requestId)
        if (userId != null) {
            MDC.put("user_id", userId)
        }
        MDC.put("request_path", path)
        MDC.put("request_method", method)
        
        logger.debug { "Incoming request: $method $path" }
        
        val startTime = System.currentTimeMillis()
        
        return chain.filter(exchange)
            .doFinally {
                val duration = System.currentTimeMillis() - startTime
                val statusCode = exchange.response.statusCode?.value() ?: 0
                
                logger.info { "Request completed: $method $path - status=$statusCode duration=${duration}ms" }
                
                // Очищаем MDC
                MDC.remove("request_id")
                MDC.remove("trace_id")
                MDC.remove("user_id")
                MDC.remove("request_path")
                MDC.remove("request_method")
            }
    }
}

