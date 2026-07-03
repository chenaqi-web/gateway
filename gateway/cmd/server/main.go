package main

import (
	"log"

	"backend/gateway/internal/config"
	"backend/gateway/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	r := router.New(cfg)
	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("gateway server: %v", err)
	}
}
