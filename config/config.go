// Package config provides configuration management for loading environment variables
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment
//
// This structure contains all configuration parameters needed by the
// application, primarily API credentials and endpoints.
type Config struct {
	// OpenAI API key for authentication
	OpenAIAPIKey string

	// OpenAI API endpoint URL
	OpenAIAPIURL string
}

// Load reads configuration from .env file
//
// Reads the .env file from the current directory and parses key-value
// pairs. Currently supports OPENAI_API_KEY and OPENAI_API_URL variables.
// Lines starting with '#' are treated as comments.
func Load() (*Config, error) {
	// Try to load .env file from current directory
	envPath := filepath.Join(".", ".env")
	if err := godotenv.Load(envPath); err != nil {
		// Check if it's a file not found error
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("configuration file .env not found")
		}
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Create config struct
	config := &Config{
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
		OpenAIAPIURL: os.Getenv("OPENAI_API_URL"),
	}

	// Validate required fields
	if config.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not found in .env file")
	}

	if config.OpenAIAPIURL == "" {
		// Set default URL if not provided
		config.OpenAIAPIURL = "https://api.openai.com/v1/chat/completions"
	}

	return config, nil
}