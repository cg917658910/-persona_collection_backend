package main

import (
	"log"

	"pm-backend/internal/api/router"
	"pm-backend/internal/config"
)

func main() {
	cfg := config.Load()
	r := router.NewHTTPServer(cfg)

	addr := ":" + cfg.Port
	log.Printf("starting %s on %s", cfg.AppName, addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
