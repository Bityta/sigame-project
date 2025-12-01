package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	
	"github.com/sigame/auth/internal/domain"
)

// AuthService handles authentication and authorization operations
type AuthService struct {
	userRepo  domain.UserRepository
	cacheRepo domain.CacheRepository
	jwtService *JWTService
	rateLimitConfig RateLimitConfig
}

// RateLimitConfig contains configuration for rate limiting
type RateLimitConfig struct {
	Attempts int           // Maximum number of attempts allowed
	Window   time.Duration // Time window for rate limiting
}

// NewAuthService creates a new authentication service instance
func NewAuthService(
	userRepo domain.UserRepository,
	cacheRepo domain.CacheRepository,
	jwtService *JWTService,
	rateLimitConfig RateLimitConfig,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		cacheRepo:       cacheRepo,
		jwtService:      jwtService,
		rateLimitConfig: rateLimitConfig,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	// Normalize and validate username
	normalizedUsername := domain.NormalizeUsername(req.Username)
	if err := domain.ValidateUsername(normalizedUsername); err != nil {
		return nil, err
	}

	// Validate password
	if err := domain.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check if username exists (with caching)
	exists, err := s.checkUsernameExists(ctx, normalizedUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, domain.ErrUsernameExists
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), domain.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Username:     normalizedUsername,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store session in Redis (for active sessions tracking)
	tokenClaims, _ := s.jwtService.ValidateAccessToken(accessToken)
	if tokenClaims != nil {
		_ = s.cacheRepo.CacheSession(ctx, tokenClaims.ID, user, s.jwtService.AccessTokenTTL)
	}

	// Store refresh token in database
	if err := s.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Cache username as taken
	_ = s.cacheRepo.CacheUsernameExists(ctx, normalizedUsername, true, domain.UsernameExistsCacheTTL)

	return &domain.RegisterResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest, clientIP string) (*domain.LoginResponse, error) {
	// Check rate limit
	if err := s.checkRateLimit(ctx, clientIP); err != nil {
		return nil, err
	}

	// Normalize username
	normalizedUsername := domain.NormalizeUsername(req.Username)

	// Get user by username
	user, err := s.userRepo.GetUserByUsername(ctx, normalizedUsername)
	if err != nil {
		s.incrementRateLimit(ctx, clientIP)
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.incrementRateLimit(ctx, clientIP)
		return nil, domain.ErrInvalidCredentials
	}

	// Reset rate limit on successful login
	_ = s.cacheRepo.ResetRateLimit(ctx, clientIP)

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store session in Redis (for active sessions tracking)
	tokenClaims, _ := s.jwtService.ValidateAccessToken(accessToken)
	if tokenClaims != nil {
		_ = s.cacheRepo.CacheSession(ctx, tokenClaims.ID, user, s.jwtService.AccessTokenTTL)
	}

	// Store refresh token in database
	if err := s.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout logs out a user by blacklisting their access token
func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	// Extract claims from token
	claims, err := s.jwtService.ExtractClaims(accessToken)
	if err != nil {
		return domain.ErrInvalidToken
	}

	// Add token to blacklist
	tokenID := domain.GetTokenID(claims)
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		if err := s.cacheRepo.AddToBlacklist(ctx, tokenID, ttl); err != nil {
			return fmt.Errorf("failed to blacklist token: %w", err)
		}
	}

	// Delete all user's refresh tokens
	if err := s.userRepo.DeleteUserRefreshTokens(ctx, claims.UserID); err != nil {
		return fmt.Errorf("failed to delete refresh tokens: %w", err)
	}

	return nil
}

// CheckUsernameAvailability checks if a username is available
func (s *AuthService) CheckUsernameAvailability(ctx context.Context, username string) (*domain.CheckUsernameResponse, error) {
	// Normalize username
	normalizedUsername := domain.NormalizeUsername(username)
	
	// Validate format
	if err := domain.ValidateUsername(normalizedUsername); err != nil {
		return &domain.CheckUsernameResponse{
			Available: false,
			Reason:    "Username format is invalid. Must be 5-50 characters, only letters, digits, underscore and hyphen.",
		}, nil
	}

	// Check if exists (with caching)
	exists, err := s.checkUsernameExists(ctx, normalizedUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}

	if exists {
		return &domain.CheckUsernameResponse{
			Available: false,
			Reason:    "Username is already taken.",
		}, nil
	}

	return &domain.CheckUsernameResponse{
		Available: true,
	}, nil
}

// GetUserInfo retrieves user information by ID
func (s *AuthService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// RefreshTokens refreshes access and refresh tokens
func (s *AuthService) RefreshTokens(ctx context.Context, refreshTokenString string) (*domain.TokenPair, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Check if token is blacklisted
	tokenID := domain.GetTokenID(claims)
	blacklisted, err := s.cacheRepo.IsBlacklisted(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to check blacklist: %w", err)
	}
	if blacklisted {
		return nil, domain.ErrTokenBlacklisted
	}

	// Verify token exists in database
	tokenHash := domain.HashToken(refreshTokenString)
	storedToken, err := s.userRepo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	// Check if token is expired
	if storedToken.IsExpired() {
		_ = s.userRepo.DeleteRefreshToken(ctx, tokenHash)
		return nil, domain.ErrTokenExpired
	}

	// Get user
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	newAccessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Delete old refresh token and store new one
	_ = s.userRepo.DeleteRefreshToken(ctx, tokenHash)
	if err := s.storeRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// ValidateAccessToken validates an access token and checks blacklist
func (s *AuthService) ValidateAccessToken(ctx context.Context, tokenString string) (*domain.Claims, error) {
	// Validate token signature and expiration
	claims, err := s.jwtService.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check if token is blacklisted
	tokenID := domain.GetTokenID(claims)
	blacklisted, err := s.cacheRepo.IsBlacklisted(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to check blacklist: %w", err)
	}
	if blacklisted {
		return nil, domain.ErrTokenBlacklisted
	}

	return claims, nil
}

// Helper methods

func (s *AuthService) checkUsernameExists(ctx context.Context, username string) (bool, error) {
	// Check cache first
	exists, found, err := s.cacheRepo.GetUsernameExists(ctx, username)
	if err != nil {
		return false, err
	}
	if found {
		return exists, nil
	}

	// Check database
	exists, err = s.userRepo.UsernameExists(ctx, username)
	if err != nil {
		return false, err
	}

	// Cache result
	_ = s.cacheRepo.CacheUsernameExists(ctx, username, exists, domain.UsernameExistsCacheTTL)

	return exists, nil
}

func (s *AuthService) storeRefreshToken(ctx context.Context, userID uuid.UUID, tokenString string) error {
	tokenHash := domain.HashToken(tokenString)
	refreshToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTokenTTL()),
	}

	return s.userRepo.CreateRefreshToken(ctx, refreshToken)
}

func (s *AuthService) checkRateLimit(ctx context.Context, ip string) error {
	count, err := s.cacheRepo.GetRateLimitCount(ctx, ip)
	if err != nil {
		return err
	}

	if count >= s.rateLimitConfig.Attempts {
		return domain.ErrRateLimitExceeded
	}

	return nil
}

func (s *AuthService) incrementRateLimit(ctx context.Context, ip string) {
	_, _ = s.cacheRepo.IncrementRateLimit(ctx, ip, s.rateLimitConfig.Window)
}

