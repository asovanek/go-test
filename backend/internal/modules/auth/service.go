package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/example/authapp/backend/internal/modules/user"
	"github.com/example/authapp/backend/internal/platform/authn"
	"github.com/example/authapp/backend/internal/platform/config"
	"github.com/example/authapp/backend/internal/platform/events"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service orchestrates authentication operations.
type Service struct {
	repo *user.Repository
	bus  *events.Bus
	cfg  *config.Config
}

func NewService(repo *user.Repository, bus *events.Bus, cfg *config.Config) *Service {
	return &Service{repo: repo, bus: bus, cfg: cfg}
}

var (
	errEmailTaken        = errors.New("email already registered")
	errInvalidCredential = errors.New("invalid credentials")
)

// SignUp creates a user and returns a JWT.
func (s *Service) SignUp(ctx context.Context, email, password string) (*AuthTokensResponse, error) {
	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	id := uuid.New()
	now := time.Now().UTC()
	u := &user.User{
		ID:           id,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	s.bus.Publish(ctx, UserSignedUp{
		UserID:    id,
		Email:     email,
		CreatedAt: now,
	})

	token, err := authn.MintToken(s.cfg, id, email)
	if err != nil {
		return nil, err
	}
	return &AuthTokensResponse{
		Token: token,
		User: AuthUserCompact{
			ID:    id.String(),
			Email: email,
		},
	}, nil
}

// SignIn validates credentials and returns a JWT.
func (s *Service) SignIn(ctx context.Context, email, password string) (*AuthTokensResponse, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errInvalidCredential
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errInvalidCredential
	}
	token, err := authn.MintToken(s.cfg, u.ID, u.Email)
	if err != nil {
		return nil, err
	}
	_ = ctx // reserved for auditing hooks
	return &AuthTokensResponse{
		Token: token,
		User: AuthUserCompact{
			ID:    u.ID.String(),
			Email: u.Email,
		},
	}, nil
}

// HTTP helpers for translating service errors.

func signupHTTPStatus(err error) int {
	if errors.Is(err, errEmailTaken) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func signinHTTPStatus(err error) int {
	if errors.Is(err, errInvalidCredential) {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}
