package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	QueryParamUserID   = "user_id"
	QueryParamToken   = "token"
	ErrorInvalidGameID = "Invalid game ID"
	ErrorUserIDRequired = "user_id is required"
	ErrorInvalidUserID  = "Invalid user ID"
	ErrorTokenRequired  = "token is required"
	ErrorInvalidToken   = "Invalid or expired token"
	ErrorGameNotFound   = "Game not found or not started"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

