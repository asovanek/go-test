package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/authapp/backend/internal/platform/authn"
	"github.com/example/authapp/backend/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.OpenSQLiteDB(t)
	cfg := testutil.TestConfig(t)

	engine := gin.New()
	Register(engine.Group("/api/v1"), db, cfg)

	t.Run("unauthorized without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d", rec.Code)
		}
	})

	t.Run("returns profile with valid token", func(t *testing.T) {
		id := uuid.New()
		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
		if err != nil {
			t.Fatalf("hash password: %v", err)
		}
		now := time.Now().UTC()
		repo := NewRepository(db)
		if err := repo.Create(&User{
			ID:           id,
			Email:        "me@example.com",
			PasswordHash: string(hash),
			CreatedAt:    now,
			UpdatedAt:    now,
		}); err != nil {
			t.Fatalf("create user: %v", err)
		}

		token, err := authn.MintToken(cfg, id, "me@example.com")
		if err != nil {
			t.Fatalf("mint token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
		}

		var profile MeResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &profile); err != nil {
			t.Fatalf("decode profile: %v", err)
		}
		if profile.Email != "me@example.com" {
			t.Fatalf("email = %q", profile.Email)
		}
		if profile.ID != id.String() {
			t.Fatalf("id = %q, want %q", profile.ID, id.String())
		}
	})
}
