// Main entry point for the mushroom classifier application
package main

import (
	"log"

	"github.com/mushroom-classifier/mushroom-classifier-go/config"
	"github.com/mushroom-classifier/mushroom-classifier-go/gui"
)

func main() {
	// Load configuration from .env file
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and setup GUI
	app, err := gui.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create GUI: %v", err)
	}

	// Run the application
	app.Run()
}