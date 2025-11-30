package com.sigame.lobby.security

import java.util.UUID

@Target(AnnotationTarget.VALUE_PARAMETER)
@Retention(AnnotationRetention.RUNTIME)
annotation class CurrentUser

data class AuthenticatedUser(
    val userId: UUID,
    val username: String
)

