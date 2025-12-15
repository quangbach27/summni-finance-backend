package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
)

type DatabaseConfig struct {
	// DATABASE CONFIG
	host            string
	database        string
	user            string
	password        string
	port            string
	maxConns        int32
	minConns        int32
	maxConnLifeTime int32
	maxConnIdleTime int32
}

// Database getters
func (db DatabaseConfig) Host() string           { return db.host }
func (db DatabaseConfig) Name() string           { return db.database }
func (db DatabaseConfig) User() string           { return db.user }
func (db DatabaseConfig) Password() string       { return db.password }
func (db DatabaseConfig) Port() string           { return db.port }
func (db DatabaseConfig) MaxConns() int32        { return db.maxConns }
func (db DatabaseConfig) MinConns() int32        { return db.minConns }
func (db DatabaseConfig) MaxConnLifeTime() int32 { return db.maxConnLifeTime }
func (db DatabaseConfig) MaxConnIdleTime() int32 { return db.maxConnIdleTime }

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

func loadConfig() *Config {
	return &Config{
		database: DatabaseConfig{
			host:            getEnv("POSTGRES_HOST", "localhost"),
			database:        getEnv("POSTGRES_DATABASE", "db"),
			user:            getEnv("POSTGRES_USER", "user"),
			password:        getEnv("POSTGRES_PASSWORD", "password"),
			port:            getEnv("POSTGRES_PORT", "5432"),
			maxConns:        getEnvAsInt32("MAX_CONNS", 5),
			minConns:        getEnvAsInt32("MIN_CONNS", 1),
			maxConnLifeTime: getEnvAsInt32("MAX_CONN_LIFE_TIME", 30),
			maxConnIdleTime: getEnvAsInt32("MAX_CONN_IDLE_TIME", 30),
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

func getEnvAsInt32(key string, defaultValue int32) int32 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	formatedValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		slog.Warn(fmt.Sprintf("Could not parse environment variable %s='%s' as integer: %s. Using default value %d.", key, value, err.Error(), defaultValue))
	}

	return int32(formatedValue)
}
