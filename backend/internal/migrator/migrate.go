package migrator

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run applies SQL migrations from migrationsDir (absolute or relative path to .sql files).
func Run(databaseURL string, migrationsDir string) error {
	dir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("migrations dir: %w", err)
	}
	sourceURL := "file://" + filepath.ToSlash(dir)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("migrate new: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
