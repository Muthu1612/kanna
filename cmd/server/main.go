package main

import (
	"log"

	"github.com/Muthu1612/kanna/internal/api"
)

func main() {
	router := api.NewRouter()

	log.Println("Starting server on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
