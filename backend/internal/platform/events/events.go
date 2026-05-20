package events

import "context"

// Event represents a typed domain event dispatched on the bus.
type Event interface {
	Type() string
}

// Subscriber handles an event asynchronously.
type Subscriber func(ctx context.Context, event Event) error
