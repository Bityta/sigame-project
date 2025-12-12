package websocket

import "time"

const (
	QueryParamUserID = "user_id"
)

const (
	ErrorInvalidGameID           = "Invalid game ID"
	ErrorInvalidUserID            = "Invalid user ID"
	ErrorUserIDRequired           = "user_id is required"
	ErrorGameNotFound             = "Game not found or not started"
	ErrorInvalidMessageFormat     = "Invalid message format"
	ErrorGameManagerNotFound      = "Game not found or not started"
	ErrorMissingServerTime        = "Missing server_time in PONG"
	ErrorInvalidServerTimeType    = "Invalid server_time type"
	ErrorInvalidRTT               = "Invalid RTT"
)

const (
	WriteWait = 10 * time.Second

	PongWait = 60 * time.Second

	JSONPingPeriod = 5 * time.Second

	MaxMessageSize = 8192

	MaxRTTSamples = 10

	MaxRTTDuration = 10 * time.Second
)

