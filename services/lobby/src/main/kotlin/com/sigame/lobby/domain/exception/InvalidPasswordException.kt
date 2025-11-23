package com.sigame.lobby.domain.exception

/**
 * Исключение выбрасывается при неверном пароле комнаты
 */
class InvalidPasswordException : RuntimeException("Invalid room password")

