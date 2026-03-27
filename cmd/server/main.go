package main

import (
	"log"

	"pm-backend/internal/api/router"
	"pm-backend/internal/config"
)

var version = "dev"

func main() {
	cfg := config.Load()
	r := router.NewHTTPServer(cfg)

	addr := ":" + cfg.Port
	log.Printf("starting %s version=%s on %s", cfg.AppName, version, addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
