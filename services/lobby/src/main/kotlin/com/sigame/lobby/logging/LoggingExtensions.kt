package com.sigame.lobby.logging

import mu.KLogger
import org.slf4j.MDC

/**
 * Extension функции для улучшения логирования
 */

/**
 * Логирует с дополнительным контекстом
 */
fun KLogger.infoWithContext(context: Map<String, String>, msg: () -> Any?) {
    val previousValues = context.mapValues { (key, _) -> MDC.get(key) }
    try {
        context.forEach { (key, value) -> MDC.put(key, value) }
        this.info(msg)
    } finally {
        previousValues.forEach { (key, value) ->
            if (value != null) MDC.put(key, value) else MDC.remove(key)
        }
    }
}

/**
 * Логирует ошибку с дополнительным контекстом
 */
fun KLogger.errorWithContext(context: Map<String, String>, throwable: Throwable? = null, msg: () -> Any?) {
    val previousValues = context.mapValues { (key, _) -> MDC.get(key) }
    try {
        context.forEach { (key, value) -> MDC.put(key, value) }
        if (throwable != null) {
            this.error(throwable, msg)
        } else {
            this.error(msg)
        }
    } finally {
        previousValues.forEach { (key, value) ->
            if (value != null) MDC.put(key, value) else MDC.remove(key)
        }
    }
}

/**
 * Выполняет блок кода с установленным MDC контекстом
 */
inline fun <T> withMDCContext(context: Map<String, String>, block: () -> T): T {
    val previousValues = context.mapValues { (key, _) -> MDC.get(key) }
    return try {
        context.forEach { (key, value) -> MDC.put(key, value) }
        block()
    } finally {
        previousValues.forEach { (key, value) ->
            if (value != null) MDC.put(key, value) else MDC.remove(key)
        }
    }
}

