package com.sigame.lobby.domain.exception

import java.util.UUID

/**
 * Исключение выбрасывается, когда комната заполнена
 */
class RoomFullException(roomId: UUID) : RuntimeException("Room $roomId is full")

