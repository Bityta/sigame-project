package com.sigame.lobby.domain.enums

/**
 * Статусы комнаты в лобби
 */
enum class RoomStatus {
    /** Комната ожидает игроков */
    WAITING,
    
    /** Комната запускается */
    STARTING,
    
    /** Игра идет */
    PLAYING,
    
    /** Игра завершена */
    FINISHED,
    
    /** Комната отменена */
    CANCELLED
}

