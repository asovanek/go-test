package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/example/authapp/backend/internal/modules/auth"
	"github.com/example/authapp/backend/internal/modules/user"
	"github.com/example/authapp/backend/internal/platform/config"
	"github.com/example/authapp/backend/internal/platform/events"
	platformmw "github.com/example/authapp/backend/internal/platform/middleware"
	"github.com/example/authapp/backend/internal/platform/validator"
)

// New builds the HTTP engine with all routes.
func New(cfg *config.Config, db *gorm.DB, bus *events.Bus, val *validator.V) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(platformmw.CORS(platformmw.CORSOptions{
		AllowedOrigins: cfg.CORSOrigins,
		AllowLocalhost: cfg.CORSAllowLocalhost,
	}))

	api := engine.Group("/api/v1")
	auth.Register(api.Group("/auth"), db, bus, cfg, val)
	user.Register(api, db, cfg)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	engine.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})

	return engine
}
