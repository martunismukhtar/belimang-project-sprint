package routes

import (
	"belimang/src/api/handlers"
	"belimang/src/pkg/merchant"

	"github.com/gofiber/fiber/v2"
)

// ActivityRouter sets up the activity routes
func MerchantRouter(app fiber.Router, merchantService merchant.Service) {

	// Create activity group with JWT middleware
	merchantGroup := app.Group("admin/merchants")

	// Activity routes
	merchantGroup.Post("/", handlers.CreateMerchant(merchantService))
	merchantGroup.Post("/:merchantId/items", handlers.CreateMerchantItems(merchantService))
	// merchantGroup.Get("/nearby/:lat/:lon", handlers.FindNearbyMerchant(purchaseService))

}
