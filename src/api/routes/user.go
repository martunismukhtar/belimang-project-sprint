package routes

import (
	"belimang/src/api/handlers"
	"belimang/src/api/middleware"
	"belimang/src/pkg/user"

	"github.com/gofiber/fiber/v2"
)

// UserRouter sets up user authentication routes
func UserRouter(app fiber.Router, userService user.Service) {
	users := app.Group("/users")

	// Public routes
	users.Post("/register", handlers.RegisterUser(userService))
	users.Post("/login", handlers.LoginUser(userService))

	// Protected routes - require JWT and user role
	users.Get("/me", middleware.JWTAuth(userService), middleware.IsUser(), handlers.GetCurrentUser(userService))
}

// AdminRouter sets up admin authentication routes
func AdminRouter(app fiber.Router, userService user.Service) {
	admin := app.Group("/admin")

	// Public routes
	admin.Post("/register", handlers.RegisterAdmin(userService))
	admin.Post("/login", handlers.LoginAdmin(userService))

	// Protected routes - require JWT and admin role
	admin.Get("/me", middleware.JWTAuth(userService), middleware.IsAdmin(), handlers.GetCurrentUser(userService))
}
