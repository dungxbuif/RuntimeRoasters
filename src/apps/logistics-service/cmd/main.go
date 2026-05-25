package main

import (
	"log"

	"RuntimeRoasters/apps/logistics-service/internal/app"
)

func main() {
	application, cleanup, err := app.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	defer cleanup()

	application.Run()
}
