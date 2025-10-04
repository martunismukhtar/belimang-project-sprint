package entities

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EstimateRequest struct {
	UserLocation struct {
		Lat  float64 `json:"lat" validate:"required"`
		Long float64 `json:"long" validate:"required"`
	} `json:"userLocation" validate:"required"`

	Orders []struct {
		MerchantID      string `json:"merchantId" validate:"required"`
		IsStartingPoint bool   `json:"isStartingPoint"`
		Items           []struct {
			ItemID   string `json:"itemId" validate:"required"`
			Quantity int    `json:"quantity" validate:"required,gt=0"`
		} `json:"items" validate:"required,dive"`
	} `json:"orders" validate:"required,min=1"`
}

type OrderRequest struct {
	CalculatedEstimateId string `json:"calculatedEstimateId" gorm:"not null" validate:"required"`
}

type OrderWrapper struct {
	Items           []OrderItemWrapper `json:"items"`
	MerchantID      uuid.UUID          `json:"merchantId"`
	IsStartingPoint bool               `json:"isStartingPoint"`
}

type OrderItemWrapper struct {
	ItemID   uuid.UUID `json:"itemId"`
	Quantity int       `json:"quantity"`
}

type OrderItem struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	MerchantID uuid.UUID `json:"merchantId" gorm:"column:merchant_id;not null" validate:"required"`
	OrderID    uuid.UUID `json:"orderId" gorm:"column:order_id;not null" validate:"required"`
	ItemID     uuid.UUID `json:"itemId" gorm:"column:item_id;not null" validate:"required"`
	Quantity   int       `json:"quantity" gorm:"column:quantity;not null;default:1" validate:"required,gt=0,number"`
	// Relasi ke Order
	Order Order `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order,omitempty"`
}

type Order struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	TotalPrice float64   `json:"totalPrice" gorm:"column:total_price;not null" validate:"required,gt=0"`
	UserID     uuid.UUID `json:"userId" gorm:"column:user_id;not null"`
	CreateAt   int64     `json:"createAt" gorm:"column:created_at;not null"`
	// Relasi ke OrderItem
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"orderItems,omitempty"`
}

type DeliveryEstimate struct {
	ID                uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey"`
	UserID            uuid.UUID       `json:"userId" gorm:"column:user_id;not null"`
	Orders            json.RawMessage `json:"orders" gorm:"column:orders;not null" validate:"required"`
	TotalPrice        float64         `json:"totalPrice" gorm:"column:total_price;not null" validate:"required,gt=0"`
	EstimatedDelivery float64         `json:"estimatedDelivery" gorm:"column:estimated_delivery_time_minutes;not null" validate:"required"`
}

func (u *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.ID = id
	}
	return nil
}

func (u *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.ID = id
	}
	return nil
}

func (u *DeliveryEstimate) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.ID = id
	}
	return nil
}

func (m DeliveryEstimate) TableName() string {
	return "delivery_estimate"
}

func (m Order) TableName() string {
	return "orders"
}

func (m OrderItem) TableName() string {
	return "order_items"
}
