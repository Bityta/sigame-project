package com.sigame.lobby.controller

import com.sigame.lobby.domain.dto.*
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.exception.ErrorResponse
import com.sigame.lobby.security.AuthenticatedUser
import com.sigame.lobby.security.CurrentUser
import com.sigame.lobby.service.domain.RoomLifecycleService
import com.sigame.lobby.service.domain.RoomMembershipService
import com.sigame.lobby.service.domain.RoomQueryService
import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.Parameter
import io.swagger.v3.oas.annotations.media.Content
import io.swagger.v3.oas.annotations.media.Schema
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import mu.KotlinLogging
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*
import java.util.UUID

private val logger = KotlinLogging.logger {}

/**
 * REST контроллер для управления лобби комнат
 */
@RestController
@RequestMapping("/api/lobby")
@Tag(name = "Lobby", description = "API для управления игровыми комнатами и лобби")
class LobbyController(
    private val roomQueryService: RoomQueryService,
    private val roomLifecycleService: RoomLifecycleService,
    private val roomMembershipService: RoomMembershipService
) {
    
    /**
     * Health check endpoint
     */
    @GetMapping("/health")
    @Operation(
        summary = "Проверка здоровья сервиса",
        description = "Возвращает статус работоспособности сервиса"
    )
    suspend fun health(): ResponseEntity<Map<String, String>> {
        return ResponseEntity.ok(mapOf("status" to "UP"))
    }
    
    /**
     * Создание новой комнаты
     */
    @PostMapping("/rooms")
    @Operation(
        summary = "Создать новую игровую комнату",
        description = "Создает новую игровую комнату с указанными параметрами"
    )
    @ApiResponses(
        value = [
            ApiResponse(
                responseCode = "201",
                description = "Комната успешно создана",
                content = [Content(schema = Schema(implementation = RoomDto::class))]
            ),
            ApiResponse(
                responseCode = "400",
                description = "Некорректные параметры запроса",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            ),
            ApiResponse(
                responseCode = "401",
                description = "Пользователь не авторизован"
            )
        ]
    )
    suspend fun createRoom(
        @Valid @RequestBody request: CreateRoomRequest,
        @Parameter(hidden = true) @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<RoomDto> {
        logger.info { 
            "User ${user.userId} creating room: name='${request.name}', " +
            "packId=${request.packId}, maxPlayers=${request.maxPlayers}, " +
            "isPublic=${request.isPublic}, hasPassword=${request.password != null}, " +
            "settings=${request.settings}"
        }
        val room = roomLifecycleService.createRoom(user.userId, request)
        return ResponseEntity.status(HttpStatus.CREATED).body(room)
    }
    
    /**
     * Получение списка комнат с фильтрацией и пагинацией
     */
    @GetMapping("/rooms")
    @Operation(
        summary = "Получить список комнат",
        description = "Возвращает список публичных комнат с пагинацией и фильтрацией"
    )
    @ApiResponses(
        value = [
            ApiResponse(
                responseCode = "200",
                description = "Список комнат получен",
                content = [Content(schema = Schema(implementation = RoomListResponse::class))]
            )
        ]
    )
    suspend fun getRooms(
        @Parameter(description = "Номер страницы (начиная с 0)")
        @RequestParam(defaultValue = "0") page: Int,
        @Parameter(description = "Размер страницы")
        @RequestParam(defaultValue = "20") size: Int,
        @Parameter(description = "Фильтр по статусу комнаты")
        @RequestParam(required = false) status: String?,
        @Parameter(description = "Только комнаты со свободными местами")
        @RequestParam(name = "has_slots", required = false) hasSlots: Boolean?
    ): ResponseEntity<RoomListResponse> {
        val statusEnum = status?.let {
            try {
                RoomStatus.valueOf(it.uppercase())
            } catch (e: IllegalArgumentException) {
                logger.warn { "Invalid status value: $it" }
                null
            }
        }
        
        val response = roomQueryService.getRooms(page, size, statusEnum, hasSlots)
        return ResponseEntity.ok(response)
    }
    
    /**
     * Получение комнаты по ID
     */
    @GetMapping("/rooms/{id}")
    @Operation(
        summary = "Получить комнату по ID",
        description = "Возвращает детальную информацию о комнате"
    )
    @ApiResponses(
        value = [
            ApiResponse(
                responseCode = "200",
                description = "Комната найдена",
                content = [Content(schema = Schema(implementation = RoomDto::class))]
            ),
            ApiResponse(
                responseCode = "404",
                description = "Комната не найдена",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            )
        ]
    )
    suspend fun getRoomById(
        @Parameter(description = "ID комнаты")
        @PathVariable id: UUID
    ): ResponseEntity<RoomDto> {
        val room = roomQueryService.getRoomById(id)
        return ResponseEntity.ok(room)
    }
    
    /**
     * Получение комнаты по коду
     */
    @GetMapping("/rooms/code/{code}")
    @Operation(
        summary = "Получить комнату по коду",
        description = "Возвращает детальную информацию о комнате по её коду"
    )
    @ApiResponses(
        value = [
            ApiResponse(
                responseCode = "200",
                description = "Комната найдена",
                content = [Content(schema = Schema(implementation = RoomDto::class))]
            ),
            ApiResponse(
                responseCode = "404",
                description = "Комната не найдена",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            )
        ]
    )
    suspend fun getRoomByCode(
        @Parameter(description = "Код комнаты (6 символов)")
        @PathVariable code: String
    ): ResponseEntity<RoomDto> {
        val room = roomQueryService.getRoomByCode(code)
        return ResponseEntity.ok(room)
    }
    
    /**
     * Присоединение к комнате
     */
    @PostMapping("/rooms/{id}/join")
    @Operation(
        summary = "Присоединиться к комнате",
        description = "Добавляет текущего пользователя в комнату"
    )
    @ApiResponses(
        value = [
            ApiResponse(
                responseCode = "200",
                description = "Успешно присоединились",
                content = [Content(schema = Schema(implementation = RoomDto::class))]
            ),
            ApiResponse(
                responseCode = "400",
                description = "Неверный пароль или некорректные параметры",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            ),
            ApiResponse(
                responseCode = "404",
                description = "Комната не найдена",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            ),
            ApiResponse(
                responseCode = "409",
                description = "Комната заполнена или игрок уже в комнате",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            )
        ]
    )
    suspend fun joinRoom(
        @Parameter(description = "ID комнаты")
        @PathVariable id: UUID,
        @RequestBody(required = false) request: JoinRoomRequest?,
        @Parameter(hidden = true) @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<RoomDto> {
        logger.info { "User ${user.userId} joining room $id" }
        val joinRequest = request ?: JoinRoomRequest()
        val room = roomMembershipService.joinRoom(id, user.userId, joinRequest)
        return ResponseEntity.ok(room)
    }
    
    /**
     * Выход из комнаты
     */
    @DeleteMapping("/rooms/{id}/leave")
    @Operation(
        summary = "Покинуть комнату",
        description = "Удаляет текущего пользователя из комнаты"
    )
    @ApiResponses(
        value = [
            ApiResponse(
                responseCode = "204",
                description = "Успешно покинули комнату"
            ),
            ApiResponse(
                responseCode = "404",
                description = "Комната не найдена или игрок не в комнате",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            )
        ]
    )
    suspend fun leaveRoom(
        @Parameter(description = "ID комнаты")
        @PathVariable id: UUID,
        @Parameter(hidden = true) @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<Void> {
        logger.info { "User ${user.userId} leaving room $id" }
        roomMembershipService.leaveRoom(id, user.userId)
        return ResponseEntity.noContent().build()
    }
    
    /**
     * Старт игры в комнате
     */
    @PostMapping("/rooms/{id}/start")
    @Operation(
        summary = "Запустить игру",
        description = "Запускает игру в комнате (только для хоста)"
    )
    @ApiResponses(
        value = [
            ApiResponse(
                responseCode = "200",
                description = "Игра успешно запущена",
                content = [Content(schema = Schema(implementation = StartGameResponse::class))]
            ),
            ApiResponse(
                responseCode = "400",
                description = "Недостаточно игроков для старта",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            ),
            ApiResponse(
                responseCode = "403",
                description = "Только хост может запустить игру",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            ),
            ApiResponse(
                responseCode = "404",
                description = "Комната не найдена",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            ),
            ApiResponse(
                responseCode = "409",
                description = "Комната не в состоянии ожидания",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            )
        ]
    )
    suspend fun startRoom(
        @Parameter(description = "ID комнаты")
        @PathVariable id: UUID,
        @Parameter(hidden = true) @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<StartGameResponse> {
        logger.info { "User ${user.userId} starting room $id" }
        val response = roomLifecycleService.startRoom(id, user.userId)
        return ResponseEntity.ok(response)
    }
    
    /**
     * Обновление настроек комнаты
     */
    @PatchMapping("/rooms/{id}/settings")
    @Operation(
        summary = "Обновить настройки комнаты",
        description = "Обновляет настройки комнаты (только для хоста)"
    )
    @ApiResponses(
        value = [
            ApiResponse(
                responseCode = "200",
                description = "Настройки обновлены",
                content = [Content(schema = Schema(implementation = RoomDto::class))]
            ),
            ApiResponse(
                responseCode = "403",
                description = "Только хост может обновлять настройки",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            ),
            ApiResponse(
                responseCode = "404",
                description = "Комната не найдена",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            ),
            ApiResponse(
                responseCode = "409",
                description = "Нельзя изменить настройки после старта игры",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            )
        ]
    )
    suspend fun updateRoomSettings(
        @Parameter(description = "ID комнаты")
        @PathVariable id: UUID,
        @Valid @RequestBody request: UpdateRoomSettingsRequest,
        @Parameter(hidden = true) @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<RoomDto> {
        logger.info { "User ${user.userId} updating room $id settings" }
        val room = roomLifecycleService.updateRoomSettings(id, user.userId, request)
        return ResponseEntity.ok(room)
    }
    
    /**
     * Удаление комнаты
     */
    @DeleteMapping("/rooms/{id}")
    @Operation(
        summary = "Удалить комнату",
        description = "Удаляет комнату (только для хоста)"
    )
    @ApiResponses(
        value = [
            ApiResponse(
                responseCode = "204",
                description = "Комната удалена"
            ),
            ApiResponse(
                responseCode = "403",
                description = "Только хост может удалить комнату",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            ),
            ApiResponse(
                responseCode = "404",
                description = "Комната не найдена",
                content = [Content(schema = Schema(implementation = ErrorResponse::class))]
            )
        ]
    )
    suspend fun deleteRoom(
        @Parameter(description = "ID комнаты")
        @PathVariable id: UUID,
        @Parameter(hidden = true) @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<Void> {
        logger.info { "User ${user.userId} deleting room $id" }
        roomLifecycleService.deleteRoom(id, user.userId)
        return ResponseEntity.noContent().build()
    }
}
