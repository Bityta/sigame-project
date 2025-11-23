package domain

import (
	"errors"
	"regexp"
)

var (
	// ErrUserNotFound is returned when user is not found
	ErrUserNotFound = errors.New("user not found")
	
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	
	// ErrUsernameExists is returned when username already exists
	ErrUsernameExists = errors.New("username already exists")
	
	// ErrInvalidUsername is returned when username format is invalid
	ErrInvalidUsername = errors.New("invalid username format")
	
	// ErrInvalidPassword is returned when password is too weak
	ErrInvalidPassword = errors.New("invalid password")
	
	// ErrInvalidToken is returned when token is invalid
	ErrInvalidToken = errors.New("invalid token")
	
	// ErrTokenExpired is returned when token has expired
	ErrTokenExpired = errors.New("token has expired")
	
	// ErrTokenBlacklisted is returned when token is blacklisted
	ErrTokenBlacklisted = errors.New("token has been revoked")
	
	// ErrRateLimitExceeded is returned when rate limit is exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// UsernameRegex is the regex pattern for valid usernames
// Allows: latin letters, cyrillic letters, digits, underscore, hyphen
// Length: 5-50 characters
var UsernameRegex = regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9_-]{5,50}$`)

// ValidateUsername validates username format
func ValidateUsername(username string) error {
	if !UsernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}

// ValidatePassword validates password requirements
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrInvalidPassword
	}
	if len(password) > MaxPasswordLength {
		return ErrInvalidPassword
	}
	return nil
}

