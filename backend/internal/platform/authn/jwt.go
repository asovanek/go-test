package authn

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"authapp/internal/platform/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// ContextUserIDKey is the Gin context key for the authenticated user's ID (string UUID).
	ContextUserIDKey = "auth_user_id"
)

// Claims embeds standard claims with email.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// JWTMiddleware validates Bearer tokens and stores user id in context.
func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("Authorization")
		if hdr == "" || !strings.HasPrefix(strings.ToLower(hdr), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		tokenStr := strings.TrimSpace(hdr[7:])
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		sub := claims.Subject
		if sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if _, err := uuid.Parse(sub); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set(ContextUserIDKey, sub)
		c.Next()
	}
}

// MintToken issues a signed JWT for the user.
func MintToken(cfg *config.Config, userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	exp := now.Add(cfg.JWTExpiry)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(cfg.JWTSecret))
}
