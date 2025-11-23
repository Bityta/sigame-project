package com.sigame.lobby.security

import com.sigame.lobby.grpc.AuthServiceClient
import kotlinx.coroutines.reactor.mono
import mu.KotlinLogging
import org.slf4j.MDC
import org.springframework.core.Ordered
import org.springframework.core.annotation.Order
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpStatus
import org.springframework.stereotype.Component
import org.springframework.web.server.ServerWebExchange
import org.springframework.web.server.WebFilter
import org.springframework.web.server.WebFilterChain
import reactor.core.publisher.Mono
import java.util.UUID

private val logger = KotlinLogging.logger {}

/**
 * Фильтр для аутентификации пользователей через JWT токены
 * Использует Auth Service для валидации токенов
 * Order = 0 (выполняется после CorsWebFilter с HIGHEST_PRECEDENCE)
 */
@Component
@Order(0)
class AuthenticationFilter(
    private val authServiceClient: AuthServiceClient
) : WebFilter {
    
    companion object {
        private val PUBLIC_PATHS = setOf(
            "/api/lobby/health",
            "/actuator/health",
            "/actuator/prometheus",
            "/actuator/metrics",
            "/actuator/info",
            "/metrics",
            "/api-docs",
            "/swagger-ui",
            "/swagger-ui.html",
            "/v3/api-docs"
        )
        
        private val PUBLIC_METHOD_PATHS = mapOf(
            "GET" to setOf("/api/lobby/rooms")
        )
    }
    
    override fun filter(exchange: ServerWebExchange, chain: WebFilterChain): Mono<Void> {
        val path = exchange.request.path.value()
        val method = exchange.request.method.name()
        
        // Always allow OPTIONS requests (CORS preflight)
        if (method == "OPTIONS") {
            return chain.filter(exchange)
        }
        
        // Skip authentication for public paths
        if (PUBLIC_PATHS.any { path.startsWith(it) }) {
            return chain.filter(exchange)
        }
        
        // Skip authentication for specific method + path combinations
        PUBLIC_METHOD_PATHS[method]?.let { paths ->
            if (paths.any { path.startsWith(it) }) {
                return chain.filter(exchange)
            }
        }
        
        // Validate JWT token
        val authHeader = exchange.request.headers.getFirst(HttpHeaders.AUTHORIZATION)
        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            logger.warn { "Missing or invalid Authorization header for $method $path" }
            return unauthorized(exchange)
        }
        
        val token = authHeader.substring(7)
        
        return mono {
            val userInfo = authServiceClient.validateToken(token)
            if (userInfo != null) {
                logger.debug { "User ${userInfo.userId} authenticated for $method $path" }
                
                // Add user info to exchange attributes
                exchange.attributes["userId"] = userInfo.userId
                exchange.attributes["username"] = userInfo.username
                
                // Add to MDC for logging
                MDC.put("user_id", userInfo.userId.toString())
                MDC.put("username", userInfo.username)
                
                // Add X-User-ID header for downstream processing
                val mutatedExchange = exchange.mutate()
                    .request { builder ->
                        builder.header("X-User-ID", userInfo.userId.toString())
                        builder.header("X-Username", userInfo.username)
                    }
                    .build()
                
                chain.filter(mutatedExchange)
            } else {
                logger.warn { "Token validation failed for $method $path" }
                unauthorized(exchange)
            }
        }.flatMap { it }
    }
    
    /**
     * Возвращает 401 Unauthorized ответ
     * CORS заголовки добавляются автоматически через CorsWebFilter
     */
    private fun unauthorized(exchange: ServerWebExchange): Mono<Void> {
        exchange.response.statusCode = HttpStatus.UNAUTHORIZED
        return exchange.response.setComplete()
    }
}

