package routes

import (
	"belimang/src/api/handlers"
	"belimang/src/api/middleware"
	"belimang/src/pkg/image"
	"belimang/src/pkg/user"

	"github.com/gofiber/fiber/v2"
)

// ImageRouter sets up image upload routes
func ImageRouter(app fiber.Router, userService user.Service, imageService image.Service) {
	// POST /image - protected route requiring admin authentication
	app.Post("/image",
		middleware.JWTAuth(userService),
		middleware.IsAdmin(),
		handlers.UploadImage(imageService),
	)
}