package com.sigame.lobby.exception

import com.fasterxml.jackson.databind.exc.InvalidFormatException
import com.fasterxml.jackson.databind.exc.MismatchedInputException
import com.sigame.lobby.domain.exception.*
import mu.KotlinLogging
import org.springframework.core.codec.DecodingException
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.validation.FieldError
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.RestControllerAdvice
import org.springframework.web.bind.support.WebExchangeBindException
import org.springframework.web.server.ServerWebInputException
import java.time.LocalDateTime

private val logger = KotlinLogging.logger {}

data class ErrorResponse(
    val timestamp: LocalDateTime = LocalDateTime.now(),
    val status: Int,
    val error: String,
    val message: String,
    val path: String? = null,
    val details: Map<String, String>? = null
)

@RestControllerAdvice
class GlobalExceptionHandler {
    
    @ExceptionHandler(ServerWebInputException::class)
    fun handleServerWebInputException(ex: ServerWebInputException): ResponseEntity<ErrorResponse> {
        val cause = ex.cause
        val details = mutableMapOf<String, String>()
        
        when (cause) {
            is InvalidFormatException -> {
                val fieldName = cause.path.joinToString(".") { it.fieldName ?: "unknown" }
                details["field"] = fieldName
                details["value"] = cause.value?.toString() ?: "null"
                details["expectedType"] = cause.targetType.simpleName
                logger.warn { "JSON deserialization error: Invalid format for field '$fieldName', value='${cause.value}', expected type=${cause.targetType.simpleName}" }
            }
            is MismatchedInputException -> {
                val fieldName = cause.path.joinToString(".") { it.fieldName ?: "unknown" }
                details["field"] = fieldName
                details["expectedType"] = cause.targetType?.simpleName ?: "unknown"
                logger.warn { "JSON deserialization error: Missing or mismatched field '$fieldName', expected type=${cause.targetType?.simpleName}" }
            }
            is DecodingException -> {
                val decodingError = cause as DecodingException
                logger.warn(decodingError) { "JSON decoding error: ${decodingError.message}" }
                details["error"] = decodingError.message ?: "Failed to decode JSON"
            }
            else -> {
                logger.warn(ex) { "ServerWebInputException: ${ex.message}, cause: ${cause?.message}" }
                details["error"] = ex.message ?: "Invalid input"
            }
        }
        
        val errorResponse = ErrorResponse(
            status = HttpStatus.BAD_REQUEST.value(),
            error = "Invalid Request Body",
            message = "Failed to read HTTP message: ${details["error"] ?: ex.reason ?: "Invalid JSON format"}",
            details = details
        )
        
        return ResponseEntity.badRequest().body(errorResponse)
    }
    
    @ExceptionHandler(WebExchangeBindException::class)
    fun handleValidationException(ex: WebExchangeBindException): ResponseEntity<ErrorResponse> {
        val errors = ex.bindingResult.allErrors.associate { error ->
            val fieldName = (error as? FieldError)?.field ?: "unknown"
            val errorMessage = error.defaultMessage ?: "Invalid value"
            fieldName to errorMessage
        }
        
        val errorResponse = ErrorResponse(
            status = HttpStatus.BAD_REQUEST.value(),
            error = "Validation Failed",
            message = "Invalid input parameters",
            details = errors
        )
        
        logger.warn { "Validation error: $errors" }
        return ResponseEntity.badRequest().body(errorResponse)
    }
    
    @ExceptionHandler(
        RoomNotFoundException::class,
        RoomNotFoundByCodeException::class
    )
    fun handleNotFoundException(ex: RuntimeException): ResponseEntity<ErrorResponse> {
        val errorResponse = ErrorResponse(
            status = HttpStatus.NOT_FOUND.value(),
            error = "Not Found",
            message = ex.message ?: "Resource not found"
        )
        
        logger.warn { "Not found: ${ex.message}" }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(errorResponse)
    }
    
    @ExceptionHandler(
        RoomFullException::class,
        PlayerAlreadyInRoomException::class,
        InvalidRoomStateException::class
    )
    fun handleConflictException(ex: RuntimeException): ResponseEntity<ErrorResponse> {
        val errorResponse = ErrorResponse(
            status = HttpStatus.CONFLICT.value(),
            error = "Conflict",
            message = ex.message ?: "Operation not allowed due to conflict"
        )
        
        logger.warn { "Conflict: ${ex.message}" }
        return ResponseEntity.status(HttpStatus.CONFLICT).body(errorResponse)
    }
    
    @ExceptionHandler(
        UnauthorizedRoomActionException::class
    )
    fun handleUnauthorizedException(ex: UnauthorizedRoomActionException): ResponseEntity<ErrorResponse> {
        val errorResponse = ErrorResponse(
            status = HttpStatus.FORBIDDEN.value(),
            error = "Forbidden",
            message = ex.message ?: "You do not have permission to perform this action"
        )
        
        logger.warn { "Forbidden: ${ex.message}" }
        return ResponseEntity.status(HttpStatus.FORBIDDEN).body(errorResponse)
    }
    
    @ExceptionHandler(
        InvalidPasswordException::class,
        PlayerNotInRoomException::class,
        PackNotFoundException::class,
        PackNotApprovedException::class,
        PackNotOwnedException::class,
        InsufficientPlayersException::class
    )
    fun handleBadRequestException(ex: RuntimeException): ResponseEntity<ErrorResponse> {
        val errorResponse = ErrorResponse(
            status = HttpStatus.BAD_REQUEST.value(),
            error = "Bad Request",
            message = ex.message ?: "Invalid request"
        )
        
        logger.warn { "Bad request: ${ex.message}" }
        return ResponseEntity.badRequest().body(errorResponse)
    }
    
    @ExceptionHandler(IllegalArgumentException::class)
    fun handleIllegalArgumentException(ex: IllegalArgumentException): ResponseEntity<ErrorResponse> {
        val errorResponse = ErrorResponse(
            status = HttpStatus.BAD_REQUEST.value(),
            error = "Bad Request",
            message = ex.message ?: "Invalid request"
        )
        
        logger.warn { "IllegalArgumentException: ${ex.message}" }
        return ResponseEntity.badRequest().body(errorResponse)
    }
    
    @ExceptionHandler(IllegalStateException::class)
    fun handleIllegalStateException(ex: IllegalStateException): ResponseEntity<ErrorResponse> {
        val errorResponse = ErrorResponse(
            status = HttpStatus.CONFLICT.value(),
            error = "Conflict",
            message = ex.message ?: "Operation not allowed"
        )
        
        logger.warn { "IllegalStateException: ${ex.message}" }
        return ResponseEntity.status(HttpStatus.CONFLICT).body(errorResponse)
    }
    
    @ExceptionHandler(NoSuchElementException::class)
    fun handleNoSuchElementException(ex: NoSuchElementException): ResponseEntity<ErrorResponse> {
        val errorResponse = ErrorResponse(
            status = HttpStatus.NOT_FOUND.value(),
            error = "Not Found",
            message = ex.message ?: "Resource not found"
        )
        
        logger.warn { "NoSuchElementException: ${ex.message}" }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(errorResponse)
    }
    
    @ExceptionHandler(UserInfoNotFoundException::class)
    fun handleUserInfoNotFoundException(ex: UserInfoNotFoundException): ResponseEntity<ErrorResponse> {
        val errorResponse = ErrorResponse(
            status = HttpStatus.INTERNAL_SERVER_ERROR.value(),
            error = "Internal Server Error",
            message = "Failed to retrieve user information"
        )
        
        logger.error { "UserInfoNotFoundException: ${ex.message}" }
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(errorResponse)
    }
    
    @ExceptionHandler(Exception::class)
    fun handleGenericException(ex: Exception): ResponseEntity<ErrorResponse> {
        val errorResponse = ErrorResponse(
            status = HttpStatus.INTERNAL_SERVER_ERROR.value(),
            error = "Internal Server Error",
            message = "An unexpected error occurred"
        )
        
        logger.error(ex) { "Unexpected error: ${ex.message}" }
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(errorResponse)
    }
}

