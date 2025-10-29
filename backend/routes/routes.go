package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vayura/backend/internal/handler"
	"github.com/vayura/backend/pkg"
)

// SetupRoutes configures all API routes with dependency injection
func SetupRoutes(router *gin.Engine, authHandler *handler.AuthHandler, userHandler *handler.UserHandler) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Public auth routes
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		// Protected routes
		protected := api.Group("/")
		protected.Use(pkg.AuthMiddleware())
		{
			// User profile CRUD
			protected.GET("/user/profile", userHandler.GetProfile)
			protected.PUT("/user/profile", userHandler.UpdateProfile)
			protected.DELETE("/user/profile", userHandler.DeleteProfile)
			protected.POST("/user/avatar", userHandler.UploadAvatar)
		}
	}
}
