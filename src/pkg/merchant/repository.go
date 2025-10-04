package merchant

import (
	"belimang/src/pkg/entities"

	"gorm.io/gorm"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	CreateMerchant(merchant *entities.Merchant) (*entities.Merchant, error)
	CreateItems(items *entities.Items) (*entities.Items, error)
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
func (r *repository) CreateMerchant(merchant *entities.Merchant) (*entities.Merchant, error) {
	if err := r.DB.Create(merchant).Error; err != nil {
		return nil, err
	}
	return merchant, nil
}

func (r *repository) CreateItems(items *entities.Items) (*entities.Items, error) {
	if err := r.DB.Create(items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) FindMerchantById(merchantId string) (*entities.Merchant, error) {
	var merchant entities.Merchant
	if err := r.DB.Where("id = ?", merchantId).First(&merchant).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}
