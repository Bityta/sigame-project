package com.sigame.lobby.controller

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.databind.SerializationFeature
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule
import com.fasterxml.jackson.module.kotlin.KotlinModule
import com.ninjasquad.springmockk.MockkBean
import com.sigame.lobby.config.WebFluxConfig
import com.sigame.lobby.domain.dto.*
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.exception.*
import com.sigame.lobby.grpc.AuthServiceClient
import com.sigame.lobby.grpc.UserInfo
import com.sigame.lobby.metrics.HttpMetrics
import com.sigame.lobby.metrics.HttpMetricsFilter
import com.sigame.lobby.security.CurrentUserArgumentResolver
import com.sigame.lobby.service.domain.RoomLifecycleService
import com.sigame.lobby.service.domain.RoomMembershipService
import com.sigame.lobby.service.domain.RoomQueryService
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.every
import io.mockk.just
import io.mockk.Runs
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.reactive.WebFluxTest
import org.springframework.context.annotation.Import
import org.springframework.http.MediaType
import org.springframework.http.codec.json.Jackson2JsonDecoder
import org.springframework.http.codec.json.Jackson2JsonEncoder
import org.springframework.test.context.ActiveProfiles
import org.springframework.test.web.reactive.server.WebTestClient
import org.springframework.web.reactive.function.client.ExchangeStrategies
import java.time.LocalDateTime
import java.util.UUID

/**
 * Интеграционные тесты для LobbyController
 * Тестирует все REST API эндпоинты согласно README спецификации
 */
@WebFluxTest(controllers = [LobbyController::class])
@Import(
    com.sigame.lobby.config.JacksonConfig::class,
    WebFluxConfig::class, 
    CurrentUserArgumentResolver::class, 
    HttpMetricsFilter::class,
    com.sigame.lobby.security.AuthenticationFilter::class,
    com.sigame.lobby.exception.GlobalExceptionHandler::class
)
@ActiveProfiles("test")
class LobbyControllerIntegrationTest {

    @Autowired
    private lateinit var webTestClient: WebTestClient

    @MockkBean
    private lateinit var roomQueryService: RoomQueryService

    @MockkBean
    private lateinit var roomLifecycleService: RoomLifecycleService

    @MockkBean
    private lateinit var roomMembershipService: RoomMembershipService

    @MockkBean
    private lateinit var authServiceClient: AuthServiceClient

    @MockkBean
    private lateinit var httpMetrics: HttpMetrics

    // Тестовые данные
    private val testUserId = UUID.randomUUID()
    private val testUsername = "testPlayer"
    private val testRoomId = UUID.randomUUID()
    private val testPackId = UUID.randomUUID()
    private val testRoomCode = "ABC123"
    
    // Токен для аутентифицированных запросов
    private val testToken = "Bearer test-token"
    
    // ObjectMapper с поддержкой Kotlin
    private val objectMapper = ObjectMapper().apply {
        registerModule(KotlinModule.Builder().build())
        registerModule(JavaTimeModule())
        configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false)
    }

    @BeforeEach
    fun setUp() {
        // Мокаем аутентификацию - токен test-token всегда валиден
        coEvery { authServiceClient.validateToken("test-token") } returns UserInfo(
            userId = testUserId,
            username = testUsername,
            avatarUrl = null
        )
        
        // Мокаем метрики
        every { httpMetrics.recordHttpRequest(any(), any(), any(), any()) } just Runs
    }
    
    // Утилита для сериализации запросов
    private fun toJson(obj: Any): String = objectMapper.writeValueAsString(obj)

    // ═══════════════════════════════════════════════════════════════════════════
    // GET /api/lobby/health
    // ═══════════════════════════════════════════════════════════════════════════
    
    @Nested
    @DisplayName("GET /api/lobby/health")
    inner class HealthEndpoint {

        @Test
        fun `должен вернуть статус UP`() {
            webTestClient.get()
                .uri("/api/lobby/health")
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.status").isEqualTo("UP")
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // POST /api/lobby/rooms - Создание комнаты
    // ═══════════════════════════════════════════════════════════════════════════

    @Nested
    @DisplayName("POST /api/lobby/rooms - Создание комнаты")
    inner class CreateRoomEndpoint {

        @Test
        fun `должен создать комнату с валидными данными и вернуть 201`() {
            val request = CreateRoomRequest(
                name = "Моя игра",
                packId = testPackId,
                maxPlayers = 6,
                isPublic = true,
                settings = RoomSettingsDto(
                    timeForAnswer = 30,
                    timeForChoice = 60,
                    allowWrongAnswer = true,
                    showRightAnswer = true
                )
            )

            val expectedRoom = createTestRoomDto()

            coEvery { roomLifecycleService.createRoom(testUserId, any()) } returns expectedRoom

            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isCreated
                .expectBody()
                .jsonPath("$.id").isNotEmpty
                .jsonPath("$.roomCode").isEqualTo(testRoomCode)
                .jsonPath("$.name").isEqualTo("Моя игра")
                .jsonPath("$.status").isEqualTo("WAITING")
                .jsonPath("$.maxPlayers").isEqualTo(6)
                .jsonPath("$.isPublic").isEqualTo(true)

            coVerify(exactly = 1) { roomLifecycleService.createRoom(testUserId, any()) }
        }

        @Test
        fun `должен вернуть 400 если name пустое`() {
            val request = mapOf(
                "name" to "",
                "packId" to testPackId.toString()
            )

            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `должен вернуть 400 если packId отсутствует`() {
            val request = mapOf(
                "name" to "Моя игра"
            )

            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `должен вернуть 400 если maxPlayers меньше 2`() {
            val request = mapOf(
                "name" to "Моя игра",
                "packId" to testPackId.toString(),
                "maxPlayers" to 1
            )

            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `должен вернуть 400 если maxPlayers больше 12`() {
            val request = mapOf(
                "name" to "Моя игра",
                "packId" to testPackId.toString(),
                "maxPlayers" to 15
            )

            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `должен вернуть 401 без авторизации`() {
            val request = CreateRoomRequest(
                name = "Моя игра",
                packId = testPackId
            )

            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(request)
                .exchange()
                .expectStatus().isUnauthorized
        }

        @Test
        fun `должен создать приватную комнату с паролем`() {
            val request = CreateRoomRequest(
                name = "Приватная игра",
                packId = testPackId,
                maxPlayers = 4,
                isPublic = false,
                password = "secret123"
            )

            val expectedRoom = createTestRoomDto().copy(
                isPublic = false,
                hasPassword = true
            )

            coEvery { roomLifecycleService.createRoom(testUserId, any()) } returns expectedRoom

            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(toJson(request))
                .exchange()
                .expectStatus().isCreated
                .expectBody()
                .jsonPath("$.isPublic").isEqualTo(false)
                .jsonPath("$.hasPassword").isEqualTo(true)
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // GET /api/lobby/rooms - Список комнат
    // ═══════════════════════════════════════════════════════════════════════════

    @Nested
    @DisplayName("GET /api/lobby/rooms - Список комнат")
    inner class GetRoomsEndpoint {

        @Test
        fun `должен вернуть список комнат с пагинацией`() {
            val rooms = listOf(createTestRoomDto(), createTestRoomDto())
            val response = RoomListResponse(
                rooms = rooms,
                page = 0,
                size = 20,
                totalElements = 2,
                totalPages = 1
            )

            coEvery { roomQueryService.getRooms(0, 20, null, null) } returns response

            webTestClient.get()
                .uri("/api/lobby/rooms")
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.rooms").isArray
                .jsonPath("$.rooms.length()").isEqualTo(2)
                .jsonPath("$.page").isEqualTo(0)
                .jsonPath("$.size").isEqualTo(20)
                .jsonPath("$.totalElements").isEqualTo(2)
                .jsonPath("$.totalPages").isEqualTo(1)
        }

        @Test
        fun `должен фильтровать по статусу`() {
            val response = RoomListResponse(
                rooms = listOf(createTestRoomDto()),
                page = 0,
                size = 20,
                totalElements = 1,
                totalPages = 1
            )

            coEvery { roomQueryService.getRooms(0, 20, "WAITING", null) } returns response

            webTestClient.get()
                .uri("/api/lobby/rooms?status=WAITING")
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.rooms.length()").isEqualTo(1)

            coVerify { roomQueryService.getRooms(0, 20, "WAITING", null) }
        }

        @Test
        fun `должен фильтровать комнаты со свободными местами`() {
            val response = RoomListResponse(
                rooms = listOf(createTestRoomDto()),
                page = 0,
                size = 20,
                totalElements = 1,
                totalPages = 1
            )

            coEvery { roomQueryService.getRooms(0, 20, null, true) } returns response

            webTestClient.get()
                .uri("/api/lobby/rooms?has_slots=true")
                .exchange()
                .expectStatus().isOk

            coVerify { roomQueryService.getRooms(0, 20, null, true) }
        }

        @Test
        fun `должен поддерживать кастомную пагинацию`() {
            val response = RoomListResponse(
                rooms = emptyList(),
                page = 2,
                size = 10,
                totalElements = 25,
                totalPages = 3
            )

            coEvery { roomQueryService.getRooms(2, 10, null, null) } returns response

            webTestClient.get()
                .uri("/api/lobby/rooms?page=2&size=10")
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.page").isEqualTo(2)
                .jsonPath("$.size").isEqualTo(10)
                .jsonPath("$.totalPages").isEqualTo(3)
        }

        @Test
        fun `должен вернуть пустой список если комнат нет`() {
            val response = RoomListResponse(
                rooms = emptyList(),
                page = 0,
                size = 20,
                totalElements = 0,
                totalPages = 0
            )

            coEvery { roomQueryService.getRooms(any(), any(), any(), any()) } returns response

            webTestClient.get()
                .uri("/api/lobby/rooms")
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.rooms").isArray
                .jsonPath("$.rooms.length()").isEqualTo(0)
                .jsonPath("$.totalElements").isEqualTo(0)
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // GET /api/lobby/rooms/{id} - Комната по ID
    // ═══════════════════════════════════════════════════════════════════════════

    @Nested
    @DisplayName("GET /api/lobby/rooms/{id} - Комната по ID")
    inner class GetRoomByIdEndpoint {

        @Test
        fun `должен вернуть комнату по ID`() {
            val room = createTestRoomDto()
            coEvery { roomQueryService.getRoomById(testRoomId) } returns room

            webTestClient.get()
                .uri("/api/lobby/rooms/$testRoomId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.id").isEqualTo(testRoomId.toString())
                .jsonPath("$.roomCode").isEqualTo(testRoomCode)
        }

        @Test
        fun `должен вернуть 404 если комната не найдена`() {
            val nonExistentId = UUID.randomUUID()
            coEvery { roomQueryService.getRoomById(nonExistentId) } throws RoomNotFoundException(nonExistentId)

            webTestClient.get()
                .uri("/api/lobby/rooms/$nonExistentId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNotFound
        }

        @Test
        fun `должен вернуть комнату со списком игроков`() {
            val players = listOf(
                PlayerDto(
                    userId = testUserId,
                    username = testUsername,
                    avatarUrl = null,
                    role = "HOST"
                )
            )
            val room = createTestRoomDto().copy(players = players)
            coEvery { roomQueryService.getRoomById(testRoomId) } returns room

            webTestClient.get()
                .uri("/api/lobby/rooms/$testRoomId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.players").isArray
                .jsonPath("$.players.length()").isEqualTo(1)
                .jsonPath("$.players[0].userId").isEqualTo(testUserId.toString())
                .jsonPath("$.players[0].role").isEqualTo("HOST")
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // GET /api/lobby/rooms/code/{code} - Комната по коду
    // ═══════════════════════════════════════════════════════════════════════════

    @Nested
    @DisplayName("GET /api/lobby/rooms/code/{code} - Комната по коду")
    inner class GetRoomByCodeEndpoint {

        @Test
        fun `должен вернуть комнату по коду`() {
            val room = createTestRoomDto()
            coEvery { roomQueryService.getRoomByCode(testRoomCode) } returns room

            webTestClient.get()
                .uri("/api/lobby/rooms/code/$testRoomCode")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.roomCode").isEqualTo(testRoomCode)
        }

        @Test
        fun `должен вернуть 404 если комната не найдена по коду`() {
            val invalidCode = "XXXXXX"
            coEvery { roomQueryService.getRoomByCode(invalidCode) } throws RoomNotFoundByCodeException(invalidCode)

            webTestClient.get()
                .uri("/api/lobby/rooms/code/$invalidCode")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNotFound
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // POST /api/lobby/rooms/{id}/join - Присоединиться к комнате
    // ═══════════════════════════════════════════════════════════════════════════

    @Nested
    @DisplayName("POST /api/lobby/rooms/{id}/join - Присоединиться к комнате")
    inner class JoinRoomEndpoint {

        @Test
        fun `должен присоединить пользователя к комнате`() {
            val room = createTestRoomDto().copy(currentPlayers = 2)
            coEvery { roomMembershipService.joinRoom(testRoomId, testUserId, any()) } returns room

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.currentPlayers").isEqualTo(2)

            coVerify { roomMembershipService.joinRoom(testRoomId, testUserId, any()) }
        }

        @Test
        fun `должен присоединиться с паролем к приватной комнате`() {
            val request = JoinRoomRequest(password = "secret123")
            val room = createTestRoomDto()
            coEvery { roomMembershipService.joinRoom(testRoomId, testUserId, any()) } returns room

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isOk
        }

        @Test
        fun `должен вернуть 404 если комната не найдена`() {
            coEvery { roomMembershipService.joinRoom(testRoomId, testUserId, any()) } throws 
                RoomNotFoundException(testRoomId)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNotFound
        }

        @Test
        fun `должен вернуть 409 если комната заполнена`() {
            coEvery { roomMembershipService.joinRoom(testRoomId, testUserId, any()) } throws 
                RoomFullException(testRoomId)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isEqualTo(409)
        }

        @Test
        fun `должен вернуть 409 если игрок уже в комнате`() {
            coEvery { roomMembershipService.joinRoom(testRoomId, testUserId, any()) } throws 
                PlayerAlreadyInRoomException(testUserId, testRoomId)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isEqualTo(409)
        }

        @Test
        fun `должен вернуть 400 если неверный пароль`() {
            coEvery { roomMembershipService.joinRoom(testRoomId, testUserId, any()) } throws 
                InvalidPasswordException()

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(JoinRoomRequest(password = "wrong"))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `должен вернуть 401 без авторизации`() {
            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .exchange()
                .expectStatus().isUnauthorized
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // DELETE /api/lobby/rooms/{id}/leave - Покинуть комнату
    // ═══════════════════════════════════════════════════════════════════════════

    @Nested
    @DisplayName("DELETE /api/lobby/rooms/{id}/leave - Покинуть комнату")
    inner class LeaveRoomEndpoint {

        @Test
        fun `должен позволить игроку покинуть комнату`() {
            coEvery { roomMembershipService.leaveRoom(testRoomId, testUserId) } returns Unit

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId/leave")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNoContent

            coVerify { roomMembershipService.leaveRoom(testRoomId, testUserId) }
        }

        @Test
        fun `должен вернуть 400 если игрок не в комнате`() {
            coEvery { roomMembershipService.leaveRoom(testRoomId, testUserId) } throws 
                PlayerNotInRoomException(testUserId, testRoomId)

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId/leave")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `должен вернуть 401 без авторизации`() {
            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId/leave")
                .exchange()
                .expectStatus().isUnauthorized
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // POST /api/lobby/rooms/{id}/start - Запустить игру
    // ═══════════════════════════════════════════════════════════════════════════

    @Nested
    @DisplayName("POST /api/lobby/rooms/{id}/start - Запустить игру")
    inner class StartRoomEndpoint {

        @Test
        fun `должен запустить игру и вернуть gameId и websocketUrl`() {
            val gameId = UUID.randomUUID().toString()
            val response = StartGameResponse(
                gameId = gameId,
                websocketUrl = "/api/game/$gameId/ws"
            )
            coEvery { roomLifecycleService.startRoom(testRoomId, testUserId) } returns response

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/start")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.gameId").isEqualTo(gameId)
                .jsonPath("$.websocketUrl").isEqualTo("/api/game/$gameId/ws")
        }

        @Test
        fun `должен вернуть 403 если не хост`() {
            coEvery { roomLifecycleService.startRoom(testRoomId, testUserId) } throws 
                UnauthorizedRoomActionException(testUserId, "start room", "Only the host can start the room")

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/start")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isForbidden
        }

        @Test
        fun `должен вернуть 400 если недостаточно игроков`() {
            coEvery { roomLifecycleService.startRoom(testRoomId, testUserId) } throws 
                InsufficientPlayersException(testRoomId, 1)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/start")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `должен вернуть 409 если комната не в состоянии WAITING`() {
            coEvery { roomLifecycleService.startRoom(testRoomId, testUserId) } throws 
                InvalidRoomStateException(testRoomId, RoomStatus.PLAYING, "start")

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/start")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isEqualTo(409)
        }

        @Test
        fun `должен вернуть 404 если комната не найдена`() {
            coEvery { roomLifecycleService.startRoom(testRoomId, testUserId) } throws 
                RoomNotFoundException(testRoomId)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/start")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNotFound
        }

        @Test
        fun `должен вернуть 401 без авторизации`() {
            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/start")
                .exchange()
                .expectStatus().isUnauthorized
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // PATCH /api/lobby/rooms/{id}/settings - Обновить настройки
    // ═══════════════════════════════════════════════════════════════════════════

    @Nested
    @DisplayName("PATCH /api/lobby/rooms/{id}/settings - Обновить настройки")
    inner class UpdateRoomSettingsEndpoint {

        @Test
        fun `должен обновить настройки и вернуть их`() {
            val request = UpdateRoomSettingsRequest(
                timeForAnswer = 45,
                timeForChoice = 90,
                allowWrongAnswer = false,
                showRightAnswer = true
            )
            val response = UpdateRoomSettingsResponse(
                settings = RoomSettingsDto(
                    timeForAnswer = 45,
                    timeForChoice = 90,
                    allowWrongAnswer = false,
                    showRightAnswer = true
                )
            )
            coEvery { roomLifecycleService.updateRoomSettings(testRoomId, testUserId, any()) } returns response

            webTestClient.patch()
                .uri("/api/lobby/rooms/$testRoomId/settings")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.settings.timeForAnswer").isEqualTo(45)
                .jsonPath("$.settings.timeForChoice").isEqualTo(90)
                .jsonPath("$.settings.allowWrongAnswer").isEqualTo(false)
                .jsonPath("$.settings.showRightAnswer").isEqualTo(true)
        }

        @Test
        fun `должен обновить только переданные настройки (PATCH семантика)`() {
            val request = UpdateRoomSettingsRequest(timeForAnswer = 60)
            val response = UpdateRoomSettingsResponse(
                settings = RoomSettingsDto(
                    timeForAnswer = 60,
                    timeForChoice = 60,  // Осталось по умолчанию
                    allowWrongAnswer = true,
                    showRightAnswer = true
                )
            )
            coEvery { roomLifecycleService.updateRoomSettings(testRoomId, testUserId, any()) } returns response

            webTestClient.patch()
                .uri("/api/lobby/rooms/$testRoomId/settings")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.settings.timeForAnswer").isEqualTo(60)
        }

        @Test
        fun `должен вернуть 400 если timeForAnswer меньше 10`() {
            val request = mapOf("timeForAnswer" to 5)

            webTestClient.patch()
                .uri("/api/lobby/rooms/$testRoomId/settings")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `должен вернуть 400 если timeForAnswer больше 120`() {
            val request = mapOf("timeForAnswer" to 150)

            webTestClient.patch()
                .uri("/api/lobby/rooms/$testRoomId/settings")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `должен вернуть 403 если не хост`() {
            val request = UpdateRoomSettingsRequest(timeForAnswer = 30)
            coEvery { roomLifecycleService.updateRoomSettings(testRoomId, testUserId, any()) } throws 
                UnauthorizedRoomActionException(testUserId, "update settings", "Only the host can update settings")

            webTestClient.patch()
                .uri("/api/lobby/rooms/$testRoomId/settings")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isForbidden
        }

        @Test
        fun `должен вернуть 409 если игра уже запущена`() {
            val request = UpdateRoomSettingsRequest(timeForAnswer = 30)
            coEvery { roomLifecycleService.updateRoomSettings(testRoomId, testUserId, any()) } throws 
                InvalidRoomStateException(testRoomId, RoomStatus.PLAYING, "update settings")

            webTestClient.patch()
                .uri("/api/lobby/rooms/$testRoomId/settings")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(request)
                .exchange()
                .expectStatus().isEqualTo(409)
        }

        @Test
        fun `должен вернуть 401 без авторизации`() {
            webTestClient.patch()
                .uri("/api/lobby/rooms/$testRoomId/settings")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(UpdateRoomSettingsRequest(timeForAnswer = 30))
                .exchange()
                .expectStatus().isUnauthorized
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // DELETE /api/lobby/rooms/{id} - Удалить комнату
    // ═══════════════════════════════════════════════════════════════════════════

    @Nested
    @DisplayName("DELETE /api/lobby/rooms/{id} - Удалить комнату")
    inner class DeleteRoomEndpoint {

        @Test
        fun `должен удалить комнату и вернуть 204`() {
            coEvery { roomLifecycleService.deleteRoom(testRoomId, testUserId) } returns Unit

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNoContent

            coVerify { roomLifecycleService.deleteRoom(testRoomId, testUserId) }
        }

        @Test
        fun `должен вернуть 403 если не хост`() {
            coEvery { roomLifecycleService.deleteRoom(testRoomId, testUserId) } throws 
                UnauthorizedRoomActionException(testUserId, "delete room", "Only the host can delete the room")

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isForbidden
        }

        @Test
        fun `должен вернуть 404 если комната не найдена`() {
            coEvery { roomLifecycleService.deleteRoom(testRoomId, testUserId) } throws 
                RoomNotFoundException(testRoomId)

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNotFound
        }

        @Test
        fun `должен вернуть 401 без авторизации`() {
            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId")
                .exchange()
                .expectStatus().isUnauthorized
        }
    }

    // ═══════════════════════════════════════════════════════════════════════════
    // Вспомогательные методы
    // ═══════════════════════════════════════════════════════════════════════════

    private fun createTestRoomDto(): RoomDto {
        return RoomDto(
            id = testRoomId,
            roomCode = testRoomCode,
            name = "Моя игра",
            hostId = testUserId,
            hostUsername = testUsername,
            packId = testPackId,
            packName = "Тестовый пак",
            status = "WAITING",
            maxPlayers = 6,
            currentPlayers = 1,
            isPublic = true,
            hasPassword = false,
            createdAt = LocalDateTime.now(),
            players = listOf(
                PlayerDto(
                    userId = testUserId,
                    username = testUsername,
                    avatarUrl = null,
                    role = "HOST"
                )
            ),
            settings = RoomSettingsDto(
                timeForAnswer = 30,
                timeForChoice = 60,
                allowWrongAnswer = true,
                showRightAnswer = true
            )
        )
    }
}

