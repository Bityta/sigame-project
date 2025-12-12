package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"sigame/game/internal/infrastructure/config"
)

func createHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.HTTPPort),
		Handler: router,
	}
}

func startHTTPServer(server *http.Server, addr string) error {
	return server.ListenAndServe()
}

