package handlers

import (
	"belimang/src/api/presenter"
	"belimang/src/pkg/book"
	"belimang/src/pkg/entities"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// validator instance

// CreateBook is handler/controller which creates Books in the BookShop
// @Summary      Create a new book
// @Description  Add a new book to the collection
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        book  body      entities.Book  true  "Book object"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/book [post]
func CreateBook(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.Book
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		if requestBody.Author == "" || requestBody.Title == "" {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.BookErrorResponse(errors.New(
				"Please specify title and author")))
		}
		result, err := service.InsertBook(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		return c.JSON(presenter.BookSuccessResponse(result))
	}
}

// UpdateBook is handler/controller which updates data of Books in the BookShop
// @Summary      Update a book
// @Description  Update an existing book (partial updates allowed)
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        book  body      entities.UpdateBookRequest  true  "Book update request (partial fields allowed)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/book [put]
func UpdateBook(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.UpdateBookRequest
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.BookErrorResponse(err))
		}

		// Convert partial update request to full Book entity
		bookToUpdate := entities.Book{
			ID: requestBody.ID,
		}

		if requestBody.Title != nil {
			bookToUpdate.Title = *requestBody.Title
		}
		if requestBody.Author != nil {
			bookToUpdate.Author = *requestBody.Author
		}

		result, err := service.UpdateBook(&bookToUpdate)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		return c.JSON(presenter.BookSuccessResponse(result))
	}
}

// DeleteBook is handler/controller which removes Books from the BookShop
// @Summary      Delete a book
// @Description  Remove a book from the collection
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        deleteRequest  body      entities.DeleteRequest  true  "Delete request with book ID"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/books [delete]
func DeleteBook(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.DeleteRequest
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		bookID := requestBody.ID
		fmt.Printf("Attempting to delete book with ID: %d\n", bookID)

		err = service.RemoveBook(bookID)
		if err != nil {
			fmt.Printf("Delete failed with error: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.BookErrorResponse(err))
		}
		fmt.Printf("Book with ID %d deleted successfully\n", bookID)
		return c.JSON(&fiber.Map{
			"status": true,
			"data":   "deleted successfully",
			"err":    nil,
		})
	}
}
