package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORS_AllowsLocalhostAnyPort(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORS(CORSOptions{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowLocalhost: true,
	}))
	r.POST("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	req.Header.Set("Origin", "http://localhost:5174")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5174" {
		t.Fatalf("Allow-Origin = %q", got)
	}
}

func TestCORS_RejectsUnknownOriginWhenLocalhostDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORS(CORSOptions{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowLocalhost: false,
	}))
	r.POST("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	req.Header.Set("Origin", "http://localhost:5174")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
}
