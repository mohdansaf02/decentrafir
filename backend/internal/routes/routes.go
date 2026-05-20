package routes

import (
	"github.com/decntrafir/backend/internal/config"
	"github.com/decntrafir/backend/internal/controllers"
	"github.com/decntrafir/backend/internal/middleware"
	"github.com/decntrafir/backend/internal/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	Auth     *controllers.AuthController
	FIR      *controllers.FIRController
	Evidence *controllers.EvidenceController
	Admin    *controllers.AdminController
}

func Setup(r *gin.Engine, cfg *config.Config, h Handlers) {
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.RateLimitMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/register", h.Auth.Register)
	r.POST("/login", h.Auth.Login)
	r.POST("/refresh", h.Auth.Refresh)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	auth := middleware.AuthMiddleware(cfg.JWTSecret)
	{
		api.GET("/me", auth, h.Auth.Me)

		fir := api.Group("/fir", auth)
		{
			fir.POST("/create", h.FIR.Create)
			fir.GET("/all", h.FIR.List)
			fir.GET("/analytics", middleware.RoleMiddleware(models.RolePolice, models.RoleAdmin), h.FIR.Analytics)
			fir.GET("/:id", h.FIR.Get)
			fir.PUT("/update/:id", middleware.RoleMiddleware(models.RolePolice, models.RoleAdmin), h.FIR.Update)
		}

		api.POST("/evidence/upload", auth, h.Evidence.Upload)

		admin := api.Group("/admin", auth, middleware.RoleMiddleware(models.RoleAdmin))
		{
			admin.GET("/users", h.Admin.ListUsers)
			admin.POST("/users", h.Admin.CreateUser)
			admin.DELETE("/users/:id", h.Admin.DeleteUser)
		}

		api.GET("/notifications", auth, h.Admin.Notifications)
	}
}
