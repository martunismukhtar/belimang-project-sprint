package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory string

const (
	Beverage   ProductCategory = "Beverage"
	Food       ProductCategory = "Food"
	Snack      ProductCategory = "Snack"
	Condiments ProductCategory = "Condiments"
	Additions  ProductCategory = "Additions"
)

type Items struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string          `json:"name" gorm:"column:name;not null" validate:"required,min=3,max=30"`
	ProductCategory ProductCategory `gorm:"column:product_category;not null" json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
	Price           float64         `json:"price" gorm:"column:price;type:numeric(10,2);not null" validate:"required,gt=0"`
	ImageUrl        string          `json:"imageUrl" gorm:"column:image_url;not null" validate:"required,url"`
	MerchantID      uuid.UUID       `json:"merchantId" gorm:"column:merchant_id;not null" validate:"required"`
	CreatedAt       int64           `gorm:"column:created_at;not null" json:"createdAt"`
	// MerchantID      uuid.UUID       `gorm:"type:uuid;not null" json:"merchantId"`
	Merchant Merchant `gorm:"foreignKey:MerchantID;references:ID" json:"-"`
}

type RequestItems struct {
	Name            string          `json:"name" gorm:"column:name;not null" validate:"required,min=3,max=30"`
	ProductCategory ProductCategory `gorm:"column:product_category;not null" json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
	Price           float64         `json:"price" gorm:"column:price;not null" validate:"required,gt=0"`
	ImageUrl        string          `json:"imageUrl" gorm:"column:image_url;not null" validate:"required,url"`
}

func (u *Items) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.ID = id
	}
	return nil
}
