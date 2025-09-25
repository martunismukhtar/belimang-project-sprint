package config

import (
	"belimang/src/api/routes"
	"belimang/src/pkg/book"

	"gorm.io/gorm"
)

// InitServices initializes all application services
func InitServices(db *gorm.DB) routes.Services {
	// Initialize repositories

	bookRepo := book.NewRepo(db)
	bookService := book.NewService(bookRepo)

	return routes.Services{
		BookService: bookService,
	}
}
