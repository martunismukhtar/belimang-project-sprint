package routes

import (
	"belimang/src/api/handlers"

	"belimang/src/api/middleware"
	"belimang/src/pkg/purchase"
	"belimang/src/pkg/user"

	"github.com/gofiber/fiber/v2"
)

// ActivityRouter sets up the activity routes
func PurchaseRouter(app fiber.Router, userService user.Service, service purchase.Service) {

	purchaseGroup := app.Group("/merchants")

	purchaseGroup.Get("/nearby/:lat/:lon", middleware.JWTAuth(userService), handlers.FindNearbyMerchant(service))
	app.Post("users/estimate", middleware.JWTAuth(userService), handlers.Estimate(service))
	app.Post("/users/orders", middleware.JWTAuth(userService), handlers.Order(service))
	app.Get("/users/orders", middleware.JWTAuth(userService), handlers.GetOrder(service))
}
