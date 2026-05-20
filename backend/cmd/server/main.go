// Package main runs the HTTP API server.
//
//	@title						Auth API
//	@version					1.0
//	@description				Signup, signin, and user profile API with JWT authentication.
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"authapp/docs"
	"authapp/internal/migrator"
	"authapp/internal/platform/config"
	"authapp/internal/platform/database"
	"authapp/internal/platform/events"
	"authapp/internal/platform/logger"
	platformvalidator "authapp/internal/platform/validator"
	"authapp/internal/router"
	"authapp/internal/subscribers"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(".env")
	if err != nil {
		return err
	}

	log, err := logger.New()
	if err != nil {
		return err
	}
	defer func() { _ = log.Sync() }()

	docs.SwaggerInfo.BasePath = "/api/v1"

	db, err := database.Open(cfg.DatabaseURL, true)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}

	migrationsDir := os.Getenv("MIGRATIONS_PATH")
	if migrationsDir == "" {
		migrationsDir = filepath.Join(".", "migrations")
	}
	if err := migrator.Run(cfg.DatabaseURL, migrationsDir); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	bus := events.NewBus(log)
	subscribers.RegisterAll(bus, log)

	val := platformvalidator.New()
	engine := router.New(cfg, db, bus, val)

	log.Info("listening", zap.String("addr", ":"+cfg.Port))
	return engine.Run(":" + cfg.Port)
}
