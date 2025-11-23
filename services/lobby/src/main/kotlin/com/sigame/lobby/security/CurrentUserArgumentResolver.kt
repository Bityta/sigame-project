package com.sigame.lobby.security

import kotlinx.coroutines.reactor.awaitSingleOrNull
import mu.KotlinLogging
import org.springframework.core.MethodParameter
import org.springframework.http.HttpStatus
import org.springframework.stereotype.Component
import org.springframework.web.reactive.BindingContext
import org.springframework.web.reactive.result.method.HandlerMethodArgumentResolver
import org.springframework.web.server.ResponseStatusException
import org.springframework.web.server.ServerWebExchange
import reactor.core.publisher.Mono
import java.util.UUID

private val logger = KotlinLogging.logger {}

/**
 * Resolver для автоматического извлечения пользователя из контекста запроса
 */
@Component
class CurrentUserArgumentResolver : HandlerMethodArgumentResolver {
    
    override fun supportsParameter(parameter: MethodParameter): Boolean {
        return parameter.hasParameterAnnotation(CurrentUser::class.java) &&
                parameter.parameterType == AuthenticatedUser::class.java
    }
    
    override fun resolveArgument(
        parameter: MethodParameter,
        bindingContext: BindingContext,
        exchange: ServerWebExchange
    ): Mono<Any> {
        return Mono.defer {
            // Извлекаем userId из атрибутов exchange (установлено AuthenticationFilter)
            val userId = exchange.attributes["userId"] as? UUID
            val username = exchange.attributes["username"] as? String
            
            if (userId != null && username != null) {
                Mono.just(AuthenticatedUser(userId, username))
            } else {
                // Попробуем извлечь из заголовков как fallback
                val userIdHeader = exchange.request.headers.getFirst("X-User-ID")
                val usernameHeader = exchange.request.headers.getFirst("X-Username")
                
                if (userIdHeader != null && usernameHeader != null) {
                    try {
                        val parsedUserId = UUID.fromString(userIdHeader)
                        Mono.just(AuthenticatedUser(parsedUserId, usernameHeader))
                    } catch (e: Exception) {
                        logger.error { "Failed to parse X-User-ID header: $userIdHeader" }
                        Mono.error(ResponseStatusException(HttpStatus.UNAUTHORIZED, "Invalid user ID"))
                    }
                } else {
                    logger.warn { "No user information found in request" }
                    Mono.error(ResponseStatusException(HttpStatus.UNAUTHORIZED, "User not authenticated"))
                }
            }
        }
    }
}

