package user

import (
	"time"

	"github.com/google/uuid"
)

// User maps to users table managed by golang-migrate.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName keeps GORM table name aligned with migrations.
func (User) TableName() string {
	return "users"
}
