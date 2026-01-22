package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the bot configuration
type Config struct {
	PhoneNumber   string
	DataDir       string
	CommandPrefix string
	LogLevel      string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		PhoneNumber:   getEnv("SIGNAL_PHONE_NUMBER", ""),
		DataDir:       getEnv("SIGNAL_DATA_DIR", "./signal-data"),
		CommandPrefix: getEnv("BOT_COMMAND_PREFIX", "!"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if config.PhoneNumber == "" {
		log.Fatal("SIGNAL_PHONE_NUMBER is required")
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
