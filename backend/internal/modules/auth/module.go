package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"authapp/internal/modules/user"
	"authapp/internal/platform/config"
	"authapp/internal/platform/events"
	"authapp/internal/platform/validator"
)

// Register mounts auth routes under /auth.
func Register(rg *gin.RouterGroup, db *gorm.DB, bus *events.Bus, cfg *config.Config, val *validator.V) {
	repo := user.NewRepository(db)
	svc := NewService(repo, bus, cfg)
	h := NewHandler(svc, val)

	rg.POST("/signup", h.SignUp)
	rg.POST("/signin", h.SignIn)
}
