package main

import (
	"log"

	"github.com/Muthu1612/kanna/internal/api"
	"github.com/Muthu1612/kanna/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	router := api.NewRouter()

	addr := cfg.Server.Host + ":" + cfg.Server.Port

	log.Printf("Kanna server starting on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
