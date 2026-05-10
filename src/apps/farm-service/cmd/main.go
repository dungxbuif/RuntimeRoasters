package main

import (
	"fmt"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/app"
)

func main() {
	fmt.Println("Starting...")
	application, _, err := app.InitializeApp()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		panic(err)
	}

	fmt.Println("Running...")
	application.Run()
}
