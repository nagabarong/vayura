package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vayura/backend/config"
	"github.com/vayura/backend/internal/handler"
	"github.com/vayura/backend/internal/models"
	"github.com/vayura/backend/internal/repository"
	"github.com/vayura/backend/internal/service"
	"github.com/vayura/backend/pkg"
	"github.com/vayura/backend/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := config.InitDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected successfully")

	// Set legacy global DB for backward compatibility
	config.DB = db
	repository.SetDB(db)

	// Initialize JWT secret
	pkg.SetJWTSecret(cfg.JWT.Secret)
	// Initialize JWT expiration
	pkg.SetJWTExpiration(cfg.JWT.ExpiresIn)

	// Run migrations
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	storageService := service.NewStorageService(cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, storageService)

	// Setup router and routes
	r := gin.Default()
	routes.SetupRoutes(r, authHandler, userHandler)

	// Start server
	port := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	if err := r.Run(port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
