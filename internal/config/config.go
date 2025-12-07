package config

import (
	"os"
	"strings"
	"sync"
)

type DatabaseConfig struct {
	// DATABASE CONFIG
	host     string
	database string
	user     string
	password string
}

// Database getters
func (db DatabaseConfig) Host() string     { return db.host }
func (db DatabaseConfig) Name() string     { return db.database }
func (db DatabaseConfig) User() string     { return db.user }
func (db DatabaseConfig) Password() string { return db.password }

// APP CONFIG
type AppConfig struct {
	port           string
	allowedOrigins []string
	env            string
}

func (a AppConfig) Port() string { return a.port }
func (a AppConfig) Env() string  { return a.env }
func (a AppConfig) AllowedOrigins() []string {
	return a.allowedOrigins
}

// Config holds all application configuration
// CONFIG ROOT
type Config struct {
	database DatabaseConfig
	app      AppConfig
}

func (c *Config) Database() DatabaseConfig { return c.database }
func (c *Config) App() AppConfig           { return c.app }

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

// INTERNAL LOADING
func loadConfig() *Config {
	return &Config{
		database: DatabaseConfig{
			host:     getEnv("POSTGRES_ADDR", "localhost"),
			database: getEnv("POSTGRES_DATABASE", "db"),
			user:     getEnv("POSTGRES_USER", "user"),
			password: getEnv("POSTGRES_PASSWORD", "password"),
		},
		app: AppConfig{
			port:           getEnv("PORT", "8080"),
			env:            getEnv("ENV", "dev"),
			allowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
