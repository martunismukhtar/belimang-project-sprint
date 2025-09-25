package routes

import (
	"belimang/src/api/handlers"
	"belimang/src/pkg/book"

	"github.com/gofiber/fiber/v2"
)

// ActivityRouter sets up the activity routes
func BookRouter(app fiber.Router, bookService book.Service) {
	// Create activity group with JWT middleware
	bookGroup := app.Group("/book")

	// Activity routes
	bookGroup.Post("/", handlers.CreateBook(bookService))
	bookGroup.Put("/:id", handlers.UpdateBook(bookService))
	bookGroup.Delete("/:id", handlers.DeleteBook(bookService))
}
