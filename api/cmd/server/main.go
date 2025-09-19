package main

import (
	"smartlock/internal/auth"
	"smartlock/internal/config"
	"smartlock/internal/httpserver"
	"smartlock/internal/mqtt"
	"smartlock/internal/state"
)

func main() {
	cfg := config.Load()

	store := state.NewStore(cfg.MaxHistory)

	mq := mqtt.New(cfg, store)
	mq.Start()

	authSvc := auth.New(cfg.JWTSecret, cfg.JWTExpiry)

	server := httpserver.NewServer(cfg, store, mq, authSvc)
	server.Start()
}
