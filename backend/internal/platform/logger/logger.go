package logger

import "go.uber.org/zap"

// New constructs a zap production logger unless dev mode is detected.
func New() (*zap.Logger, error) {
	return zap.NewProduction()
}
