// cmd/main.go
package main

import (
	"log"
	"myapp/app"
	"myapp/internal/scheduler"
)

func main() {
	engine, err := app.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	go scheduler.Start()
	engine.Run(":8080")
}
