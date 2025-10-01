package user

import (
	"belimang/src/pkg/entities"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository interface defines user data access operations
type Repository interface {
	Create(user *entities.User) error
	FindByID(id uuid.UUID) (*entities.User, error)
	FindByUsername(username string) (*entities.User, error)
	FindByEmailAndRole(email string, role string) (*entities.User, error)
	UsernameExists(username string) (bool, error)
	EmailExistsForRole(email string, role string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository instance
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Create inserts a new user into the database
func (r *repository) Create(user *entities.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindByID retrieves a user by their ID
func (r *repository) FindByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername retrieves a user by username (username is globally unique)
func (r *repository) FindByUsername(username string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmailAndRole retrieves a user by email and role (email is unique per role)
func (r *repository) FindByEmailAndRole(email string, role string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ? AND role = ?", email, role).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UsernameExists checks if a username already exists (globally unique)
func (r *repository) UsernameExists(username string) (bool, error) {
	var count int64
	if err := r.db.Model(&entities.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// EmailExistsForRole checks if an email exists for a specific role
func (r *repository) EmailExistsForRole(email string, role string) (bool, error) {
	var count int64
	if err := r.db.Model(&entities.User{}).Where("email = ? AND role = ?", email, role).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
