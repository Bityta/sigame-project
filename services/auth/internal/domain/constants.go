package domain

import "time"

const (
	// MinPasswordLength is the minimum allowed password length
	MinPasswordLength = 8
	
	// MaxPasswordLength is the maximum allowed password length
	MaxPasswordLength = 128
	
	// BcryptCost is the cost factor for bcrypt password hashing
	BcryptCost = 12
	
	// UsernameExistsCacheTTL is the duration for caching username existence checks
	UsernameExistsCacheTTL = 5 * time.Minute
)

// HTTP status code ranges for metrics
const (
	// StatusCode2xx represents successful HTTP responses (200-299)
	StatusCode2xx = "2xx"
	
	// StatusCode3xx represents redirect HTTP responses (300-399)
	StatusCode3xx = "3xx"
	
	// StatusCode4xx represents client error HTTP responses (400-499)
	StatusCode4xx = "4xx"
	
	// StatusCode5xx represents server error HTTP responses (500-599)
	StatusCode5xx = "5xx"
	
	// StatusCodeUnknown represents unknown or invalid status codes
	StatusCodeUnknown = "unknown"
)

// Redis key prefixes for different data types
const (
	// RedisKeyUsernameExists is the prefix for username existence cache keys
	RedisKeyUsernameExists = "username:exists:"
	
	// RedisKeyBlacklist is the prefix for blacklisted token cache keys
	RedisKeyBlacklist = "blacklist:"
	
	// RedisKeyRateLimit is the prefix for rate limit counter cache keys
	RedisKeyRateLimit = "rate_limit:login:"
	
	// RedisKeySession is the prefix for user session cache keys
	RedisKeySession = "session:"
)

