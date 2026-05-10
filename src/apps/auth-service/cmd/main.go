package main

import (
	"fmt"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/app"
)

func main() {
	fmt.Println("Initializing Auth Service...")
	application, _, err := app.InitializeApp()
	if err != nil {
		fmt.Printf("Error during initialization: %v\n", err)
		panic(err)
	}

	fmt.Println("Running Auth Service...")
	application.Run()
}
