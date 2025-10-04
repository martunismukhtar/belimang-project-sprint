package config

import (
	"log"
	"os"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func NewSwagger(app *fiber.App) error {
	// Check if swagger.json file exists
	if _, err := os.Stat("./docs/swagger.json"); os.IsNotExist(err) {
		log.Printf("Warning: Swagger file not found at ./docs/swagger.json. Skipping Swagger setup.")
		return nil
	}

	// Setup Swagger middleware with error recovery
	// Example showing defer behavior:
	// func example() {
	// 	fmt.Println("1")
	// 	defer fmt.Println("4") // This runs LAST
	// 	fmt.Println("2")
	// 	defer fmt.Println("3") // This runs second-to-last
	// 	return
	// 	// Output: 1, 2, 4, 3 (deferred functions run in LIFO order)
	// }
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error initializing Swagger: %v", r)
		}
	}()

	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "BeliMang API Documentation",
		CacheAge: 86400,
	}))

	log.Println("Swagger documentation initialized successfully at /swagger")
	return nil
}
