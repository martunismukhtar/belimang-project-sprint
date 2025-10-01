package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username  string         `json:"username" gorm:"type:varchar(30);not null;uniqueIndex:idx_users_username_unique"`
	Email     string         `json:"email" gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email_role_unique"`
	Role      string         `json:"role" gorm:"type:varchar(50);not null;default:'user';uniqueIndex:idx_users_email_role_unique"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index:idx_users_deleted_at"`
}

func (User) TableName() string {
	return "users"
}
