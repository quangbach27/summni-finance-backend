package config

import (
	"os"
	"strings"
	"sync"
)

// Config holds all application configuration
type Config struct {
	Port     string
	Database DatabaseConfig
	CORS     CORSConfig
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Database string
	User     string
	Password string
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowedOrigins []string
}

var (
	configInstance *Config
	once           sync.Once
)

// GetConfig returns the singleton instance of Config
// It loads the configuration from .env file on first call
func GetConfig() *Config {
	once.Do(func() {
		cfg := loadConfig()
		configInstance = cfg
	})

	return configInstance
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_ADDR", "localhost"),
			Database: getEnv("POSTGRES_DATABASE", "db"),
			User:     getEnv("POSTGRES_USER", "user"),
			Password: getEnv("POSTGRES_PASSWORD", "password"),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
	}

	return cfg
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
