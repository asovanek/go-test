package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"authapp/internal/platform/authn"
	"authapp/internal/platform/config"
)

// Register mounts user routes.
func Register(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	protected := api.Group("")
	protected.Use(authn.JWTMiddleware(cfg))
	protected.GET("/me", h.Me)
}
