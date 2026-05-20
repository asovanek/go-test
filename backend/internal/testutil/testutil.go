package testutil

import (
	"testing"
	"time"

	"authapp/internal/platform/config"
)

// TestConfig returns config suitable for tests.
func TestConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		Port:               "8080",
		DatabaseURL:        "unused",
		JWTSecret:          "test-secret-key",
		JWTExpiry:          time.Hour,
		CORSOrigins:        []string{"http://localhost:5173"},
		CORSAllowLocalhost: true,
	}
}
