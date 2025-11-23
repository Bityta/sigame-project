package com.sigame.lobby.domain.exception

import java.util.UUID

/**
 * Исключение выбрасывается, когда комната не найдена
 */
class RoomNotFoundException(roomId: UUID) : RuntimeException("Room with id $roomId not found")

/**
 * Исключение выбрасывается, когда комната не найдена по коду
 */
class RoomNotFoundByCodeException(roomCode: String) : RuntimeException("Room with code $roomCode not found")

