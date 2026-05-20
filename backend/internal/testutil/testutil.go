package testutil

import (
	"testing"
	"time"

	"authapp/internal/modules/user"
	"authapp/internal/platform/config"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestConfig returns config suitable for tests.
func TestConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		Port:        "8080",
		DatabaseURL: "unused",
		JWTSecret:   "test-secret-key",
		JWTExpiry:   time.Hour,
		CORSOrigins: []string{"http://localhost:5173"},
	}
}

// OpenSQLiteDB opens an in-memory SQLite database with users migrated.
func OpenSQLiteDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&user.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}
