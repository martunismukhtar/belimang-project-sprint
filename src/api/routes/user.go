package routes

import (
	"belimang/src/api/handlers"
	"belimang/src/pkg/user"

	"github.com/gofiber/fiber/v2"
)

// UserRouter sets up user authentication routes
func UserRouter(app fiber.Router, userService user.Service) {
	// User registration and login (public routes)
	users := app.Group("/users")
	users.Post("/register", handlers.RegisterUser(userService))
	users.Post("/login", handlers.LoginUser(userService))
}

// AdminRouter sets up admin authentication routes
func AdminRouter(app fiber.Router, userService user.Service) {
	// Admin registration and login (public routes)
	admin := app.Group("/admin")
	admin.Post("/register", handlers.RegisterAdmin(userService))
	admin.Post("/login", handlers.LoginAdmin(userService))
}
