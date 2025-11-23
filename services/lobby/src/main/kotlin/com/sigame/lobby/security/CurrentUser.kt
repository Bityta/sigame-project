package com.sigame.lobby.security

import java.util.UUID

/**
 * Аннотация для автоматической инжекции текущего пользователя в параметры контроллера
 */
@Target(AnnotationTarget.VALUE_PARAMETER)
@Retention(AnnotationRetention.RUNTIME)
annotation class CurrentUser

/**
 * Данные текущего аутентифицированного пользователя
 */
data class AuthenticatedUser(
    val userId: UUID,
    val username: String
)

