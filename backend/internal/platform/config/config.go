package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment.
type Config struct {
	Port       string
	DatabaseURL string
	JWTSecret  string
	JWTExpiry  time.Duration
	CORSOrigin string
}

func Load(paths ...string) (*Config, error) {
	path := ".env"
	if len(paths) > 0 && paths[0] != "" {
		path = paths[0]
	}
	if _, err := os.Stat(path); err == nil {
		if err := godotenv.Load(path); err != nil {
			return nil, fmt.Errorf("load env file %s: %w", path, err)
		}
	}

	jwtExpiry := getenv("JWT_EXPIRY", "24h")
	duration, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		return nil, fmt.Errorf("parse JWT_EXPIRY: %w", err)
	}

	dbURL := getenv("DATABASE_URL", "")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecret := getenv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return &Config{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		JWTExpiry:   duration,
		CORSOrigin:  getenv("CORS_ORIGIN", "http://localhost:5173"),
	}, nil
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" && def != "" {
		return def
	}
	return v
}
