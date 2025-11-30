package com.sigame.lobby.security

import org.springframework.core.MethodParameter
import org.springframework.http.HttpStatus
import org.springframework.stereotype.Component
import org.springframework.web.reactive.BindingContext
import org.springframework.web.reactive.result.method.HandlerMethodArgumentResolver
import org.springframework.web.server.ResponseStatusException
import org.springframework.web.server.ServerWebExchange
import reactor.core.publisher.Mono
import java.util.UUID

@Component
class CurrentUserArgumentResolver : HandlerMethodArgumentResolver {

    override fun supportsParameter(parameter: MethodParameter): Boolean =
        parameter.hasParameterAnnotation(CurrentUser::class.java) &&
            parameter.parameterType == AuthenticatedUser::class.java

    override fun resolveArgument(
        parameter: MethodParameter,
        bindingContext: BindingContext,
        exchange: ServerWebExchange
    ): Mono<Any> = Mono.defer {
        val userId = exchange.attributes["userId"] as? UUID
        val username = exchange.attributes["username"] as? String

        if (userId != null && username != null) {
            Mono.just(AuthenticatedUser(userId, username))
        } else {
            Mono.error(ResponseStatusException(HttpStatus.UNAUTHORIZED, "User not authenticated"))
        }
    }
}
