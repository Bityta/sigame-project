package com.sigame.lobby.logging

import org.springframework.http.HttpHeaders
import org.springframework.stereotype.Component

@Component
class LogSanitizer {

    companion object {
        private const val HIDDEN = "***HIDDEN***"
        private const val MAX_BODY_SIZE = 10_000

        private val SENSITIVE_FIELDS = setOf(
            "password", "token", "accessToken", "refreshToken", "authorization"
        )

        private val SENSITIVE_HEADERS = setOf(
            "authorization", "x-api-key", "cookie"
        )

        private val JSON_PATTERNS = SENSITIVE_FIELDS.map { field ->
            Regex(""""$field"\s*:\s*"[^"]*"""", RegexOption.IGNORE_CASE) to """"$field":"$HIDDEN""""
        }

        private val QUERY_PATTERNS = SENSITIVE_FIELDS.map { field ->
            Regex("""$field=[^&\s]*""", RegexOption.IGNORE_CASE) to "$field=$HIDDEN"
        }
    }

    fun sanitizeHeaders(headers: HttpHeaders): Map<String, String> =
        headers.toSingleValueMap()
            .filterKeys { key -> SENSITIVE_HEADERS.none { it.equals(key, ignoreCase = true) } }

    fun sanitizeQueryParams(params: Map<String, String>): Map<String, String> =
        params.mapValues { (key, value) ->
            if (SENSITIVE_FIELDS.any { it.equals(key, ignoreCase = true) }) HIDDEN else value
        }

    fun sanitizeBody(body: String): String {
        if (body.isBlank()) return body

        var result = body
        JSON_PATTERNS.forEach { (pattern, replacement) ->
            result = pattern.replace(result, replacement)
        }
        QUERY_PATTERNS.forEach { (pattern, replacement) ->
            result = pattern.replace(result, replacement)
        }
        return result
    }

    fun truncateAndSanitize(body: String): String =
        sanitizeBody(
            if (body.length > MAX_BODY_SIZE) body.take(MAX_BODY_SIZE) + "...[TRUNCATED]" else body
        )
}

