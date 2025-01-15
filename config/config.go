package config

import (
	"log"
	"os"
	"fmt"

	"github.com/joho/godotenv"
)

// Config holds the configuration values
type Config struct {
	APIKey      string
	DatabaseURL string
}

// LoadConfig loads the configuration from the .env file
func LoadConfig() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, err
	}

	// Read values from environment variables
	config := &Config{
		APIKey:      os.Getenv("API_KEY"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	// Ensure required variables are set
	if config.APIKey == "" {
		return nil, fmt.Errorf("API_KEY is not set in the environment")
	}
	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set in the environment")
	}

	log.Println("Configuration loaded from .env file successfully!")
	return config, nil
}
