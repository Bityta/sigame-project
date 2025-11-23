package com.sigame.lobby.domain.exception

import java.util.UUID

/**
 * Исключение выбрасывается когда пак вопросов не найден
 */
class PackNotFoundException(packId: UUID) : RuntimeException("Pack with id $packId not found")

