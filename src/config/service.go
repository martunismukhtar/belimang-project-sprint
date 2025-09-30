package config

import (
	"belimang/src/api/routes"
	"belimang/src/pkg/user"
	"os"

	"gorm.io/gorm"
)

// InitServices initializes all application services
func InitServices(db *gorm.DB) routes.Services {
	// Get JWT secret from environment or use default
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}

	// Initialize repositories
	userRepo := user.NewRepository(db)

	// Initialize servicesF
	userService := user.NewService(userRepo, jwtSecret)

	return routes.Services{
		UserService: userService,
	}
}
