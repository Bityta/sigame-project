package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"github.com/sigame/auth/internal/domain"
	"github.com/sigame/auth/internal/service"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	authService *service.AuthService
}

// NewHandler creates a new HTTP handler for authentication endpoints
func NewHandler(authService *service.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

// CheckUsername godoc
// @Summary Check username availability
// @Description Check if a username is available for registration
// @Tags auth
// @Param username query string true "Username to check"
// @Success 200 {object} domain.CheckUsernameResponse
// @Failure 400 {object} domain.ErrorResponse
// @Router /auth/check-username [get]
func (h *Handler) CheckUsername(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		respondBadRequest(c, "bad_request", "Username parameter is required")
		return
	}

	response, err := h.authService.CheckUsernameAvailability(c.Request.Context(), username)
	if err != nil {
		respondInternalError(c, "internal_error", "Failed to check username availability")
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration data"
// @Success 201 {object} domain.RegisterResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 409 {object} domain.ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "bad_request", "Invalid request body")
		return
	}

	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidUsername:
			respondBadRequest(c, "invalid_username", "Username must be 5-50 characters and contain only letters, digits, underscore or hyphen")
		case domain.ErrInvalidPassword:
			respondBadRequest(c, "invalid_password", "Password must be at least 8 characters long")
		case domain.ErrUsernameExists:
			respondConflict(c, "username_exists", "Username is already taken")
		default:
			respondInternalError(c, "internal_error", "Failed to register user")
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 429 {object} domain.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "bad_request", "Invalid request body")
		return
	}

	clientIP := c.ClientIP()
	response, err := h.authService.Login(c.Request.Context(), &req, clientIP)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			respondUnauthorized(c, "invalid_credentials", "Invalid username or password")
		case domain.ErrRateLimitExceeded:
			respondTooManyRequests(c, "rate_limit_exceeded", "Too many login attempts. Please try again later")
		default:
			respondInternalError(c, "internal_error", "Failed to login")
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Get new access and refresh tokens using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshRequest true "Refresh token"
// @Success 200 {object} domain.TokenPair
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "bad_request", "Invalid request body")
		return
	}

	response, err := h.authService.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken, domain.ErrTokenExpired, domain.ErrTokenBlacklisted:
			respondUnauthorized(c, "invalid_token", "Invalid or expired refresh token")
		default:
			respondInternalError(c, "internal_error", "Failed to refresh tokens")
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate user's access token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} domain.ErrorResponse
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	token, err := ExtractBearerToken(authHeader)
	if err != nil {
		var message string
		switch err {
		case ErrAuthHeaderMissing:
			message = "Authorization header is required"
		case ErrInvalidAuthFormat:
			message = "Invalid authorization header format"
		default:
			message = "Invalid authorization"
		}
		
		respondUnauthorized(c, "unauthorized", message)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		respondInternalError(c, "internal_error", "Failed to logout")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// GetMe godoc
// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.UserResponse
// @Failure 401 {object} domain.ErrorResponse
// @Router /auth/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		respondUnauthorized(c, "unauthorized", "User not authenticated")
		return
	}

	parsedUserID, err := uuid.Parse(userID.(string))
	if err != nil {
		respondBadRequest(c, "bad_request", "Invalid user ID")
		return
	}

	user, err := h.authService.GetUserInfo(c.Request.Context(), parsedUserID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			respondNotFound(c, "user_not_found", "User not found")
			return
		}
		respondInternalError(c, "internal_error", "Failed to get user information")
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// Health godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "auth-service",
	})
}

