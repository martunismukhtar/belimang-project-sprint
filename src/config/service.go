package config

import (
	"belimang/src/api/routes"
	"belimang/src/pkg/image"
	"belimang/src/pkg/merchant"
	"belimang/src/pkg/purchase"
	"belimang/src/pkg/user"
	"os"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// InitServices initializes all application services
func InitServices(db *gorm.DB, minioClient *minio.Client) routes.Services {
	// Get JWT secret from environment or use default
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}

	// Initialize repositories
	userRepo := user.NewRepository(db)

	// Initialize services
	userService := user.NewService(userRepo, jwtSecret)
	imageService := image.NewService(minioClient)

	//merchant

	merchantRepo := merchant.NewRepo(db)
	merchantService := merchant.NewService(merchantRepo)

	//purchase
	purchaseRepo := purchase.NewRepo(db)
	purchaseService := purchase.NewService(purchaseRepo)

	//

	return routes.Services{
		UserService:  userService,
		ImageService: imageService,
		// BookService:     bookService,
		MerchantService: merchantService,
		PurchaseService: purchaseService,
	}
}
