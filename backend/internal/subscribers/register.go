package subscribers

import (
	"context"

	"github.com/example/authapp/backend/internal/modules/auth"
	"github.com/example/authapp/backend/internal/platform/events"
	"go.uber.org/zap"
)

// RegisterAll wires subscribers onto the event bus.
func RegisterAll(bus *events.Bus, log *zap.Logger) {
	bus.Subscribe(events.EventUserSignedUp, auditLog(log))
	bus.Subscribe(events.EventUserSignedUp, welcomeStub(log))
}

func auditLog(log *zap.Logger) func(ctx context.Context, e events.Event) error {
	return func(ctx context.Context, e events.Event) error {
		evt, ok := e.(auth.UserSignedUp)
		if !ok {
			return nil
		}
		log.Info("user signed up",
			zap.String("event", e.Type()),
			zap.String("user_id", evt.UserID.String()),
			zap.String("email", evt.Email),
			zap.Time("created_at", evt.CreatedAt),
		)
		_ = ctx
		return nil
	}
}

func welcomeStub(log *zap.Logger) func(ctx context.Context, e events.Event) error {
	return func(ctx context.Context, e events.Event) error {
		evt, ok := e.(auth.UserSignedUp)
		if !ok {
			return nil
		}
		log.Info("welcome subscriber stub — would queue welcome email",
			zap.String("email", evt.Email),
		)
		_ = ctx
		return nil
	}
}
