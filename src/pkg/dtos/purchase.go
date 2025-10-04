package dtos

import (
	"github.com/google/uuid"
)

type LocationResponse struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type MerchantResponse struct {
	MerchantID       uuid.UUID        `json:"merchantId"`
	Name             string           `json:"name"`
	MerchantCategory string           `json:"merchantCategory"`
	ImageURL         string           `json:"imageUrl"`
	Location         LocationResponse `json:"location"`
	CreatedAt        string           `json:"createdAt"`
}

type ItemResponse struct {
	ItemID          uuid.UUID `json:"itemId"`
	Name            string    `json:"name"`
	ProductCategory string    `json:"productCategory"`
	Price           float64   `json:"price"`
	ImageURL        string    `json:"imageUrl"`
	CreatedAt       string    `json:"createdAt"`
}

type MerchantWithItems struct {
	Merchant MerchantResponse `json:"merchant"`
	Items    []ItemResponse   `json:"items"`
}

type MetaResponse struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type NearbyMerchantResponse struct {
	Data []MerchantWithItems `json:"data"`
	Meta MetaResponse        `json:"meta"`
}

type OrderDetail struct {
	OrderID           uuid.UUID
	MerchantID        uuid.UUID
	MerchantName      string
	MerchantCategory  string
	MerchantImageURL  string
	MerchantLong      float64
	MerchantLat       float64
	MerchantCreatedAt int64
	ItemID            uuid.UUID
	ItemName          string
	ProductCategory   string
	Price             float64
	Quantity          int
	ItemImageURL      string
	ItemCreatedAt     int64
}
