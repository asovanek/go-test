package authn

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"authapp/internal/platform/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret: "test-secret-key",
		JWTExpiry: time.Hour,
	}
}

func TestMintTokenAndMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := testConfig()
	id := uuid.New()

	token, err := MintToken(cfg, id, "user@example.com")
	if err != nil {
		t.Fatalf("mint token: %v", err)
	}

	r := gin.New()
	r.GET("/protected", JWTMiddleware(cfg), func(c *gin.Context) {
		v, _ := c.Get(ContextUserIDKey)
		c.String(http.StatusOK, v.(string))
	})

	t.Run("valid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
		}
		if rec.Body.String() != id.String() {
			t.Fatalf("user id = %q, want %q", rec.Body.String(), id.String())
		}
	})

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer not-a-jwt")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
	})
}
