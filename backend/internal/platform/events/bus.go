package events

import (
	"context"

	"go.uber.org/zap"
)

const (
	EventUserSignedUp = "user.signed_up"
)

// Bus is an in-process pub/sub dispatcher.
type Bus struct {
	logger   *zap.Logger
	handlers map[string][]Subscriber
}

// NewBus creates an empty event bus.
func NewBus(log *zap.Logger) *Bus {
	if log == nil {
		log = zap.NewNop()
	}
	return &Bus{
		logger:   log,
		handlers: make(map[string][]Subscriber),
	}
}

// Subscribe registers a subscriber for eventType.
func (b *Bus) Subscribe(eventType string, handler Subscriber) {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish dispatches event to subscribers asynchronously (one goroutine each).
func (b *Bus) Publish(ctx context.Context, event Event) {
	if event == nil {
		return
	}
	typ := event.Type()
	subs := b.handlers[typ]
	for _, fn := range subs {
		handler := fn
		go func() {
			defer func() {
				if r := recover(); r != nil {
					b.logger.Error("event subscriber panic", zap.Any("panic", r), zap.String("event_type", typ))
				}
			}()
			if err := handler(ctx, event); err != nil {
				b.logger.Warn("event subscriber error", zap.Error(err), zap.String("event_type", typ))
			}
		}()
	}
}
