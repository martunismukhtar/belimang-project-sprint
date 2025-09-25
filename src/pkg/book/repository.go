package book

import (
	"belimang/src/api/presenter"
	"belimang/src/pkg/entities"

	"gorm.io/gorm"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateBook(book *entities.Book) (*entities.Book, error)
	ReadBook() (*[]presenter.Book, error)
	UpdateBook(book *entities.Book) (*entities.Book, error)
	DeleteBook(ID uint) error
}
type repository struct {
	DB *gorm.DB
}

// NewRepo is the single instance repo that is being created.
func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

// CreateBook is a GORM repository that helps to create books
func (r *repository) CreateBook(book *entities.Book) (*entities.Book, error) {
	if err := r.DB.Create(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

// ReadBook is a GORM repository that helps to fetch books
func (r *repository) ReadBook() (*[]presenter.Book, error) {
	var entityBooks []entities.Book
	// GORM automatically excludes soft-deleted records when using Find with entities.Book
	if err := r.DB.Find(&entityBooks).Error; err != nil {
		return nil, err
	}

	// Convert entities.Book to presenter.Book
	var books []presenter.Book
	for _, book := range entityBooks {
		books = append(books, presenter.Book{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
		})
	}

	return &books, nil
}

// UpdateBook is a GORM repository that helps to update books
func (r *repository) UpdateBook(book *entities.Book) (*entities.Book, error) {
	if err := r.DB.Model(book).Updates(book).Error; err != nil {
		return nil, err
	}
	var updatedBook entities.Book
	if err := r.DB.First(&updatedBook, book.ID).Error; err != nil {
		return nil, err
	}
	return &updatedBook, nil
}

// DeleteBook is a GORM repository that helps to delete books
func (r *repository) DeleteBook(ID uint) error {
	if err := r.DB.Delete(&entities.Book{}, ID).Error; err != nil {
		return err
	}
	return nil
}
