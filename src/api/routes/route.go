package routes

import (
	"belimang/src/pkg/user"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, v *viper.Viper, db *gorm.DB, services Services) {
	// API v1 group
	api := app.Group("/api/v1")

	// Setup all routes
	UserRouter(api, services.UserService)
	AdminRouter(api, services.UserService)

	// --- Health check route for Kubernetes probes ---
	app.Get("/healthz", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB() // get underlying *sql.DB from GORM
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		if err := sqlDB.Ping(); err != nil { // try pinging the DB
			return c.Status(fiber.StatusInternalServerError).SendString("Database not reachable")
		}

		return c.SendStatus(fiber.StatusOK) // 200 OK if DB is fine
	})
}

// Services struct holds all service dependencies
type Services struct {
	UserService user.Service
	// ActivityService   activity.Service
	// UploadFileService userfile.Service
}
