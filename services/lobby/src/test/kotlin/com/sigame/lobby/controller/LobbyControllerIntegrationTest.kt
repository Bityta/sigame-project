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
import com.sigame.lobby.grpc.auth.AuthServiceClient
import com.sigame.lobby.grpc.auth.UserInfo
import com.sigame.lobby.metrics.HttpMetrics
import com.sigame.lobby.metrics.HttpMetricsFilter
import com.sigame.lobby.metrics.LobbyMetrics
import com.sigame.lobby.security.CurrentUserArgumentResolver
import com.sigame.lobby.sse.service.RoomEventPublisher
import com.sigame.lobby.service.domain.RoomLifecycleService
import com.sigame.lobby.service.domain.RoomMembershipService
import com.sigame.lobby.service.domain.RoomQueryService
import io.mockk.Runs
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.every
import io.mockk.just
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.reactive.WebFluxTest
import org.springframework.context.annotation.Import
import org.springframework.http.MediaType
import org.springframework.test.context.ActiveProfiles
import org.springframework.test.web.reactive.server.WebTestClient
import java.time.LocalDateTime
import java.util.UUID

@WebFluxTest(controllers = [
    RoomQueryController::class,
    RoomLifecycleController::class,
    RoomMembershipController::class
])
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

    @MockkBean(relaxed = true)
    private lateinit var asyncLogWriter: com.sigame.lobby.logging.AsyncLogWriter

    @MockkBean(relaxed = true)
    private lateinit var roomEventPublisher: RoomEventPublisher

    @MockkBean(relaxed = true)
    private lateinit var lobbyMetrics: LobbyMetrics

    private val testUserId = UUID.randomUUID()
    private val testUsername = "testPlayer"
    private val testRoomId = UUID.randomUUID()
    private val testPackId = UUID.randomUUID()
    private val testRoomCode = "ABC123"
    private val testToken = "Bearer test-token"

    private val objectMapper = ObjectMapper().apply {
        registerModule(KotlinModule.Builder().build())
        registerModule(JavaTimeModule())
        configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false)
    }

    @BeforeEach
    fun setUp() {
        coEvery { authServiceClient.validateToken("test-token") } returns UserInfo(
            userId = testUserId,
            username = testUsername,
            avatarUrl = null
        )
        every { httpMetrics.recordHttpRequest(any(), any(), any(), any()) } just Runs
    }

    private fun toJson(obj: Any): String = objectMapper.writeValueAsString(obj)

    @Nested
    @DisplayName("GET /api/lobby/health")
    inner class HealthEndpoint {

        @Test
        fun `should return UP status`() {
            webTestClient.get()
                .uri("/api/lobby/health")
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.status").isEqualTo("UP")
        }
    }

    @Nested
    @DisplayName("POST /api/lobby/rooms - Create Room")
    inner class CreateRoomEndpoint {

        @Test
        fun `should create room with valid data and return 201`() {
            val request = CreateRoomRequest(
                name = "My Game",
                packId = testPackId,
                maxPlayers = 6,
                isPublic = true,
                settings = RoomSettingsDto(30, 60)
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
                .jsonPath("$.status").isEqualTo("WAITING")
        }

        @Test
        fun `should return 400 if name is empty`() {
            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(mapOf("name" to "", "packId" to testPackId.toString()))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 400 if maxPlayers less than 2`() {
            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(mapOf("name" to "Game", "packId" to testPackId.toString(), "maxPlayers" to 1))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 400 if maxPlayers greater than 12`() {
            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(mapOf("name" to "Game", "packId" to testPackId.toString(), "maxPlayers" to 15))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 401 without authorization`() {
            webTestClient.post()
                .uri("/api/lobby/rooms")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(CreateRoomRequest(name = "Game", packId = testPackId))
                .exchange()
                .expectStatus().isUnauthorized
        }

        @Test
        fun `should create private room with password`() {
            val request = CreateRoomRequest(
                name = "Private Game",
                packId = testPackId,
                isPublic = false,
                password = "secret123"
            )
            val expectedRoom = createTestRoomDto().copy(isPublic = false, hasPassword = true)

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

    @Nested
    @DisplayName("GET /api/lobby/rooms - List Rooms")
    inner class GetRoomsEndpoint {

        @Test
        fun `should return rooms with pagination`() {
            val response = RoomListResponse(
                rooms = listOf(createTestRoomDto(), createTestRoomDto()),
                page = 0, size = 20, totalElements = 2, totalPages = 1
            )
            coEvery { roomQueryService.getRooms(0, 20, null, null) } returns response

            webTestClient.get()
                .uri("/api/lobby/rooms")
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.rooms.length()").isEqualTo(2)
                .jsonPath("$.totalElements").isEqualTo(2)
        }

        @Test
        fun `should filter by status`() {
            val response = RoomListResponse(listOf(createTestRoomDto()), 0, 20, 1, 1)
            coEvery { roomQueryService.getRooms(0, 20, "WAITING", null) } returns response

            webTestClient.get()
                .uri("/api/lobby/rooms?status=WAITING")
                .exchange()
                .expectStatus().isOk

            coVerify { roomQueryService.getRooms(0, 20, "WAITING", null) }
        }

        @Test
        fun `should filter by has_slots`() {
            val response = RoomListResponse(listOf(createTestRoomDto()), 0, 20, 1, 1)
            coEvery { roomQueryService.getRooms(0, 20, null, true) } returns response

            webTestClient.get()
                .uri("/api/lobby/rooms?has_slots=true")
                .exchange()
                .expectStatus().isOk

            coVerify { roomQueryService.getRooms(0, 20, null, true) }
        }

        @Test
        fun `should return empty list when no rooms`() {
            val response = RoomListResponse(emptyList(), 0, 20, 0, 0)
            coEvery { roomQueryService.getRooms(any(), any(), any(), any()) } returns response

            webTestClient.get()
                .uri("/api/lobby/rooms")
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.rooms.length()").isEqualTo(0)
        }
    }

    @Nested
    @DisplayName("GET /api/lobby/rooms/my - Get My Rooms")
    inner class GetMyRoomsEndpoint {

        @Test
        fun `should return user's active room`() {
            val room = createTestRoomDto()
            coEvery { roomQueryService.getMyActiveRoom(testUserId) } returns room

            webTestClient.get()
                .uri("/api/lobby/rooms/my")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.rooms.length()").isEqualTo(1)
                .jsonPath("$.rooms[0].id").isEqualTo(testRoomId.toString())
        }

        @Test
        fun `should return empty list if user has no active room`() {
            coEvery { roomQueryService.getMyActiveRoom(testUserId) } returns null

            webTestClient.get()
                .uri("/api/lobby/rooms/my")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.rooms.length()").isEqualTo(0)
        }

        @Test
        fun `should return 401 without authorization`() {
            webTestClient.get()
                .uri("/api/lobby/rooms/my")
                .exchange()
                .expectStatus().isUnauthorized
        }
    }

    @Nested
    @DisplayName("GET /api/lobby/rooms/{id} - Get Room by ID")
    inner class GetRoomByIdEndpoint {

        @Test
        fun `should return room by ID`() {
            val room = createTestRoomDto()
            coEvery { roomQueryService.getRoomById(testRoomId) } returns room

            webTestClient.get()
                .uri("/api/lobby/rooms/$testRoomId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isOk
                .expectBody()
                .jsonPath("$.id").isEqualTo(testRoomId.toString())
        }

        @Test
        fun `should return 404 if room not found`() {
            val nonExistentId = UUID.randomUUID()
            coEvery { roomQueryService.getRoomById(nonExistentId) } throws RoomNotFoundException(nonExistentId)

            webTestClient.get()
                .uri("/api/lobby/rooms/$nonExistentId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNotFound
        }
    }

    @Nested
    @DisplayName("GET /api/lobby/rooms/code/{code} - Get Room by Code")
    inner class GetRoomByCodeEndpoint {

        @Test
        fun `should return room by code`() {
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
        fun `should return 404 if room not found by code`() {
            coEvery { roomQueryService.getRoomByCode("XXXXXX") } throws RoomNotFoundByCodeException("XXXXXX")

            webTestClient.get()
                .uri("/api/lobby/rooms/code/XXXXXX")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNotFound
        }
    }

    @Nested
    @DisplayName("POST /api/lobby/rooms/{id}/join - Join Room")
    inner class JoinRoomEndpoint {

        @Test
        fun `should join room successfully`() {
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
        }

        @Test
        fun `should join with password`() {
            val room = createTestRoomDto()
            coEvery { roomMembershipService.joinRoom(testRoomId, testUserId, any()) } returns room

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(JoinRoomRequest(password = "secret123"))
                .exchange()
                .expectStatus().isOk
        }

        @Test
        fun `should return 409 if room is full`() {
            coEvery { roomMembershipService.joinRoom(testRoomId, testUserId, any()) } throws RoomFullException(testRoomId)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isEqualTo(409)
        }

        @Test
        fun `should return 400 if invalid password`() {
            coEvery { roomMembershipService.joinRoom(testRoomId, testUserId, any()) } throws InvalidPasswordException()

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(JoinRoomRequest(password = "wrong"))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 401 without authorization`() {
            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/join")
                .contentType(MediaType.APPLICATION_JSON)
                .exchange()
                .expectStatus().isUnauthorized
        }
    }

    @Nested
    @DisplayName("DELETE /api/lobby/rooms/{id}/leave - Leave Room")
    inner class LeaveRoomEndpoint {

        @Test
        fun `should leave room successfully`() {
            coEvery { roomMembershipService.leaveRoom(testRoomId, testUserId) } returns Unit

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId/leave")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNoContent
        }

        @Test
        fun `should return 400 if player not in room`() {
            coEvery { roomMembershipService.leaveRoom(testRoomId, testUserId) } throws
                PlayerNotInRoomException(testUserId, testRoomId)

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId/leave")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isBadRequest
        }
    }

    @Nested
    @DisplayName("POST /api/lobby/rooms/{id}/kick - Kick Player")
    inner class KickPlayerEndpoint {

        @Test
        fun `should kick player successfully`() {
            val targetUserId = UUID.randomUUID()
            coEvery { roomMembershipService.kickPlayer(testRoomId, testUserId, targetUserId) } returns Unit

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/kick")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(KickPlayerRequest(targetUserId))
                .exchange()
                .expectStatus().isNoContent

            coVerify { roomMembershipService.kickPlayer(testRoomId, testUserId, targetUserId) }
        }

        @Test
        fun `should return 403 if not host`() {
            val targetUserId = UUID.randomUUID()
            coEvery { roomMembershipService.kickPlayer(testRoomId, testUserId, targetUserId) } throws
                UnauthorizedRoomActionException(testUserId, "kick player", "Only the host can kick players")

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/kick")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(KickPlayerRequest(targetUserId))
                .exchange()
                .expectStatus().isForbidden
        }

        @Test
        fun `should return 400 if trying to kick self`() {
            coEvery { roomMembershipService.kickPlayer(testRoomId, testUserId, testUserId) } throws
                CannotKickSelfException(testUserId)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/kick")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(KickPlayerRequest(testUserId))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 400 if player not in room`() {
            val targetUserId = UUID.randomUUID()
            coEvery { roomMembershipService.kickPlayer(testRoomId, testUserId, targetUserId) } throws
                PlayerNotInRoomException(targetUserId, testRoomId)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/kick")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(KickPlayerRequest(targetUserId))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 409 if room not in WAITING status`() {
            val targetUserId = UUID.randomUUID()
            coEvery { roomMembershipService.kickPlayer(testRoomId, testUserId, targetUserId) } throws
                InvalidRoomStateException(testRoomId, RoomStatus.PLAYING, "kick player")

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/kick")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(KickPlayerRequest(targetUserId))
                .exchange()
                .expectStatus().isEqualTo(409)
        }

        @Test
        fun `should return 401 without authorization`() {
            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/kick")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(KickPlayerRequest(UUID.randomUUID()))
                .exchange()
                .expectStatus().isUnauthorized
        }
    }

    @Nested
    @DisplayName("POST /api/lobby/rooms/{id}/transfer-host - Transfer Host")
    inner class TransferHostEndpoint {

        @Test
        fun `should transfer host successfully`() {
            val newHostId = UUID.randomUUID()
            coEvery { roomMembershipService.transferHostManually(testRoomId, testUserId, newHostId) } returns Unit

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/transfer-host")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(TransferHostRequest(newHostId))
                .exchange()
                .expectStatus().isNoContent

            coVerify { roomMembershipService.transferHostManually(testRoomId, testUserId, newHostId) }
        }

        @Test
        fun `should return 403 if not host`() {
            val newHostId = UUID.randomUUID()
            coEvery { roomMembershipService.transferHostManually(testRoomId, testUserId, newHostId) } throws
                UnauthorizedRoomActionException(testUserId, "transfer host", "Only the host can transfer host role")

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/transfer-host")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(TransferHostRequest(newHostId))
                .exchange()
                .expectStatus().isForbidden
        }

        @Test
        fun `should return 400 if transferring to self`() {
            coEvery { roomMembershipService.transferHostManually(testRoomId, testUserId, testUserId) } throws
                IllegalArgumentException("Cannot transfer host to yourself")

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/transfer-host")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(TransferHostRequest(testUserId))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 400 if new host not in room`() {
            val newHostId = UUID.randomUUID()
            coEvery { roomMembershipService.transferHostManually(testRoomId, testUserId, newHostId) } throws
                PlayerNotInRoomException(newHostId, testRoomId)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/transfer-host")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(TransferHostRequest(newHostId))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 409 if room not in WAITING status`() {
            val newHostId = UUID.randomUUID()
            coEvery { roomMembershipService.transferHostManually(testRoomId, testUserId, newHostId) } throws
                InvalidRoomStateException(testRoomId, RoomStatus.PLAYING, "transfer host")

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/transfer-host")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(TransferHostRequest(newHostId))
                .exchange()
                .expectStatus().isEqualTo(409)
        }

        @Test
        fun `should return 401 without authorization`() {
            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/transfer-host")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(TransferHostRequest(UUID.randomUUID()))
                .exchange()
                .expectStatus().isUnauthorized
        }
    }

    @Nested
    @DisplayName("POST /api/lobby/rooms/{id}/start - Start Game")
    inner class StartRoomEndpoint {

        @Test
        fun `should start game and return gameId and websocketUrl`() {
            val gameId = UUID.randomUUID().toString()
            val response = StartGameResponse(gameId = gameId, websocketUrl = "/api/game/$gameId/ws")
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
        fun `should return 403 if not host`() {
            coEvery { roomLifecycleService.startRoom(testRoomId, testUserId) } throws
                UnauthorizedRoomActionException(testUserId, "start room", "Only the host can start the room")

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/start")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isForbidden
        }

        @Test
        fun `should return 400 if insufficient players`() {
            coEvery { roomLifecycleService.startRoom(testRoomId, testUserId) } throws
                InsufficientPlayersException(testRoomId, 1)

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/start")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 409 if room not in WAITING status`() {
            coEvery { roomLifecycleService.startRoom(testRoomId, testUserId) } throws
                InvalidRoomStateException(testRoomId, RoomStatus.PLAYING, "start")

            webTestClient.post()
                .uri("/api/lobby/rooms/$testRoomId/start")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isEqualTo(409)
        }
    }

    @Nested
    @DisplayName("PATCH /api/lobby/rooms/{id}/settings - Update Settings")
    inner class UpdateRoomSettingsEndpoint {

        @Test
        fun `should update settings successfully`() {
            val request = UpdateRoomSettingsRequest(timeForAnswer = 45, timeForChoice = 90)
            val response = UpdateRoomSettingsResponse(
                settings = RoomSettingsDto(45, 90)
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
        }

        @Test
        fun `should return 400 if timeForAnswer less than 10`() {
            webTestClient.patch()
                .uri("/api/lobby/rooms/$testRoomId/settings")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(mapOf("timeForAnswer" to 5))
                .exchange()
                .expectStatus().isBadRequest
        }

        @Test
        fun `should return 403 if not host`() {
            coEvery { roomLifecycleService.updateRoomSettings(testRoomId, testUserId, any()) } throws
                UnauthorizedRoomActionException(testUserId, "update settings", "Only the host can update settings")

            webTestClient.patch()
                .uri("/api/lobby/rooms/$testRoomId/settings")
                .contentType(MediaType.APPLICATION_JSON)
                .header("Authorization", testToken)
                .bodyValue(UpdateRoomSettingsRequest(timeForAnswer = 30))
                .exchange()
                .expectStatus().isForbidden
        }
    }

    @Nested
    @DisplayName("DELETE /api/lobby/rooms/{id} - Delete Room")
    inner class DeleteRoomEndpoint {

        @Test
        fun `should delete room and return 204`() {
            coEvery { roomLifecycleService.deleteRoom(testRoomId, testUserId) } returns Unit

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNoContent
        }

        @Test
        fun `should return 403 if not host`() {
            coEvery { roomLifecycleService.deleteRoom(testRoomId, testUserId) } throws
                UnauthorizedRoomActionException(testUserId, "delete room", "Only the host can delete the room")

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isForbidden
        }

        @Test
        fun `should return 404 if room not found`() {
            coEvery { roomLifecycleService.deleteRoom(testRoomId, testUserId) } throws
                RoomNotFoundException(testRoomId)

            webTestClient.delete()
                .uri("/api/lobby/rooms/$testRoomId")
                .header("Authorization", testToken)
                .exchange()
                .expectStatus().isNotFound
        }
    }

    private fun createTestRoomDto(): RoomDto = RoomDto(
        id = testRoomId,
        roomCode = testRoomCode,
        name = "My Game",
        hostId = testUserId,
        hostUsername = testUsername,
        packId = testPackId,
        packName = "Test Pack",
        status = "WAITING",
        maxPlayers = 6,
        currentPlayers = 1,
        isPublic = true,
        hasPassword = false,
        createdAt = LocalDateTime.now(),
        players = listOf(PlayerDto(testUserId, testUsername, null, "HOST")),
        settings = RoomSettingsDto(30, 60)
    )
}
