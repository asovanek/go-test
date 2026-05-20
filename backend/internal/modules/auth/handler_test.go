package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"authapp/internal/testdb"
	"authapp/internal/platform/events"
	platformvalidator "authapp/internal/platform/validator"
	"authapp/internal/testutil"
	"github.com/gin-gonic/gin"
)

func setupAuthRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := testdb.OpenSQLite(t)
	bus := events.NewBus(nil)
	cfg := testutil.TestConfig(t)
	val := platformvalidator.New()

	r := gin.New()
	Register(r.Group("/auth"), db, bus, cfg, val)
	return r
}

func TestHandler_SignUp_Success(t *testing.T) {
	r := setupAuthRouter(t)
	body := `{"email":"handler@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var resp AuthTokensResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Token == "" || resp.User.Email != "handler@example.com" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestHandler_SignUp_ValidationError(t *testing.T) {
	r := setupAuthRouter(t)
	body := `{"email":"bad","password":"short"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestHandler_SignUp_Conflict(t *testing.T) {
	r := setupAuthRouter(t)
	body := `{"email":"dup@example.com","password":"password123"}`
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if i == 0 && rec.Code != http.StatusCreated {
			t.Fatalf("first signup status = %d", rec.Code)
		}
		if i == 1 && rec.Code != http.StatusConflict {
			t.Fatalf("second signup status = %d, body = %s", rec.Code, rec.Body.String())
		}
	}
}

func TestHandler_SignIn_InvalidCredentials(t *testing.T) {
	r := setupAuthRouter(t)

	signupReq := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString(`{"email":"signin@example.com","password":"password123"}`))
	signupReq.Header.Set("Content-Type", "application/json")
	signupRec := httptest.NewRecorder()
	r.ServeHTTP(signupRec, signupReq)
	if signupRec.Code != http.StatusCreated {
		t.Fatalf("setup signup status = %d", signupRec.Code)
	}

	signinReq := httptest.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBufferString(`{"email":"signin@example.com","password":"wrongpass"}`))
	signinReq.Header.Set("Content-Type", "application/json")
	signinRec := httptest.NewRecorder()
	r.ServeHTTP(signinRec, signinReq)

	if signinRec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", signinRec.Code, signinRec.Body.String())
	}

	var errBody map[string]string
	if err := json.Unmarshal(signinRec.Body.Bytes(), &errBody); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if errBody["error"] != "invalid credentials" {
		t.Fatalf("error = %q", errBody["error"])
	}
}
