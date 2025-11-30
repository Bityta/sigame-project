package com.sigame.lobby.security

import com.sigame.lobby.controller.ApiRoutes
import com.sigame.lobby.grpc.auth.AuthServiceClient
import kotlinx.coroutines.reactor.mono
import mu.KotlinLogging
import org.slf4j.MDC
import org.springframework.core.annotation.Order
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpStatus
import org.springframework.stereotype.Component
import org.springframework.web.server.ServerWebExchange
import org.springframework.web.server.WebFilter
import org.springframework.web.server.WebFilterChain
import reactor.core.publisher.Mono

private val logger = KotlinLogging.logger {}

@Component
@Order(0)
class AuthenticationFilter(
    private val authServiceClient: AuthServiceClient
) : WebFilter {

    companion object {
        private val PUBLIC_PATHS = ApiRoutes.Public.PATHS
        private val PUBLIC_GET_PATHS = ApiRoutes.Public.GET_PATHS
    }

    override fun filter(exchange: ServerWebExchange, chain: WebFilterChain): Mono<Void> {
        val path = exchange.request.path.value()
        val method = exchange.request.method.name()

        if (method == "OPTIONS" || isPublicPath(path, method)) {
            return chain.filter(exchange)
        }

        val token = extractToken(exchange) ?: run {
            logger.warn { "Missing or invalid Authorization header for $method $path" }
            return unauthorized(exchange)
        }

        return mono {
            authServiceClient.validateToken(token)?.let { userInfo ->
                exchange.attributes["userId"] = userInfo.userId
                exchange.attributes["username"] = userInfo.username
                MDC.put("user_id", userInfo.userId.toString())
                MDC.put("username", userInfo.username)

                chain.filter(
                    exchange.mutate()
                        .request { it.header("X-User-ID", userInfo.userId.toString()).header("X-Username", userInfo.username) }
                        .build()
                )
            } ?: run {
                logger.warn { "Token validation failed for $method $path" }
                unauthorized(exchange)
            }
        }.flatMap { it }
    }

    private fun isPublicPath(path: String, method: String): Boolean =
        PUBLIC_PATHS.any { path.startsWith(it) } ||
            (method == "GET" && PUBLIC_GET_PATHS.any { path.startsWith(it) })

    private fun extractToken(exchange: ServerWebExchange): String? {
        val headerToken = exchange.request.headers.getFirst(HttpHeaders.AUTHORIZATION)
            ?.takeIf { it.startsWith("Bearer ") }
            ?.substring(7)
        
        if (headerToken != null) return headerToken
        
        val path = exchange.request.path.value()
        if (path.contains("/events")) {
            return exchange.request.queryParams.getFirst("token")
        }
        
        return null
    }

    private fun unauthorized(exchange: ServerWebExchange): Mono<Void> {
        exchange.response.statusCode = HttpStatus.UNAUTHORIZED
        return exchange.response.setComplete()
    }
}
