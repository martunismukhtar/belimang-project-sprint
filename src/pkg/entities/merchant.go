package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchantCategory string

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

type Merchant struct {
	ID               uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string           `gorm:"column:name;not null" json:"name" validate:"required,min=3,max=30"`
	ImageUrl         string           `gorm:"column:image_url;not null" json:"imageUrl" validate:"required,url"`
	Lat              float64          `gorm:"column:lat;not null" json:"lat" validate:"required"`
	Long             float64          `gorm:"column:long;not null" json:"long" validate:"required"`
	MerchantCategory MerchantCategory `gorm:"column:merchant_category;not null" json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
	CreatedAt        int64            `gorm:"column:created_at;not null" json:"createdAt"`
	Items            []Items          `gorm:"foreignKey:MerchantID;references:ID" json:"items"`
}

type Location struct {
	Lat  float64 `json:"lat" validate:"required"`
	Long float64 `json:"long" validate:"required"`
}

type RequestMerchant struct {
	Name             string           `gorm:"column:name;not null" json:"name" validate:"required,min=3,max=30"`
	ImageUrl         string           `gorm:"column:image_url;not null" json:"imageUrl" validate:"required,url"`
	Location         Location         `gorm:"column:name;not null" json:"location" validate:"required"`
	MerchantCategory MerchantCategory `gorm:"column:merchant_category;not null" json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
}

func (m Merchant) TableName() string {
	return "merchants"
}

func (u *Merchant) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.ID = id
	}
	return nil
}

func (m Merchant) PrimaryKey() string {
	return "id"
}
