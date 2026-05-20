package user

import (
	"github.com/google/uuid"
)

// Service exposes read operations for user profile.
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(id uuid.UUID) (*User, error) {
	return s.repo.FindByID(id)
}
