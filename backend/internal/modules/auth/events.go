package auth

import (
	"time"

	"authapp/internal/platform/events"
	"github.com/google/uuid"
)

// UserSignedUp is emitted after a user record is successfully created.
type UserSignedUp struct {
	UserID    uuid.UUID
	Email     string
	CreatedAt time.Time
}

// Type implements events.Event.
func (UserSignedUp) Type() string {
	return events.EventUserSignedUp
}

var _ events.Event = UserSignedUp{}
