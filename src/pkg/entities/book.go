package entities

import (
	"time"

	"gorm.io/gorm"
)

// Book Constructs your Book model under entities.
type Book struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string         `json:"title" gorm:"type:varchar(255);not null"`
	Author    string         `json:"author" gorm:"type:varchar(255)"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UpdateBookRequest struct is used to parse Update Requests for Books (partial updates)
type UpdateBookRequest struct {
	ID     uint    `json:"id" binding:"required"`
	Title  *string `json:"title,omitempty"`
	Author *string `json:"author,omitempty"`
}

// DeleteRequest struct is used to parse Delete Requests for Books
type DeleteRequest struct {
	ID uint `json:"id"`
}
