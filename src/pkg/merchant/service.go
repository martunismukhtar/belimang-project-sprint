package merchant

import (
	"belimang/src/pkg/entities"

	"github.com/google/uuid"
)

// Service is an interface from which our api module can access our repository of all our models
type Service interface {
	InsertMerchant(merchant *entities.Merchant) (*entities.Merchant, error)
	CreateItems(items *entities.Items, merchantId uuid.UUID) (*entities.Items, error)
}

type service struct {
	repository Repository
}

// NewService is used to create a single instance of the service
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// InsertBook is a service layer that helps insert book in BookShop
func (s *service) InsertMerchant(merchant *entities.Merchant) (*entities.Merchant, error) {
	return s.repository.CreateMerchant(merchant)
}

func (s *service) CreateItems(items *entities.Items, merchantId uuid.UUID) (*entities.Items, error) {
	//cek apakah merchantId ada di db

	return s.repository.CreateItems(items)
}

// FetchBooks is a service layer that helps fetch all books in BookShop
// func (s *service) FetchBooks() (*[]presenter.Book, error) {
// 	return s.repository.ReadBook()
// }

// // UpdateBook is a service layer that helps update books in BookShop
// func (s *service) UpdateBook(book *entities.Book) (*entities.Book, error) {
// 	return s.repository.UpdateBook(book)
// }

// // RemoveBook is a service layer that helps remove books from BookShop
// func (s *service) RemoveBook(ID uint) error {
// 	return s.repository.DeleteBook(ID)
// }
