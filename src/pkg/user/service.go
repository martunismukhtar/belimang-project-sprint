package user

import (
	"belimang/src/pkg/entities"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrUsernameExists        = errors.New("username already exists")
	ErrEmailExists           = errors.New("email already exists for this role")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidUsername       = errors.New("username must be 5-30 characters")
	ErrInvalidPassword       = errors.New("password must be 5-30 characters")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidRole           = errors.New("invalid role")
	emailRegex               = regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
)

// Service interface defines business logic operations for users
type Service interface {
	Register(username, email, password, role string) (*entities.User, error)
	Login(username, password, role string) (string, *entities.User, error)
	GetUserByID(id string) (*entities.User, error)
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
}

type service struct {
	repository Repository
	jwtSecret  string
}

// NewService creates a new user service instance
func NewService(repo Repository, jwtSecret string) Service {
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}
	return &service{
		repository: repo,
		jwtSecret:  jwtSecret,
	}
}

// Register creates a new user account with validation
func (s *service) Register(username, email, password, role string) (*entities.User, error) {
	// Validate role
	if role != entities.RoleUser && role != entities.RoleAdmin {
		return nil, ErrInvalidRole
	}

	// Validate username length
	if len(username) < 5 || len(username) > 30 {
		return nil, ErrInvalidUsername
	}

	// Validate email format
	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}

	// Validate password length (plain text before hashing)
	if len(password) < 5 || len(password) > 30 {
		return nil, ErrInvalidPassword
	}

	// Check if username already exists (globally unique)
	exists, err := s.repository.UsernameExists(username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, ErrUsernameExists
	}

	// Check if email exists for this role (email unique per role)
	emailExists, err := s.repository.EmailExistsForRole(email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if emailExists {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entities.User{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
		Role:     role,
		Password: string(hashedPassword),
	}

	if err := s.repository.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *service) Login(username, password, role string) (string, *entities.User, error) {
	// Validate role
	if role != entities.RoleUser && role != entities.RoleAdmin {
		return "", nil, ErrInvalidRole
	}

	// Find user by username
	user, err := s.repository.FindByUsername(username)
	if err != nil {
		return "", nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}

	// Check if user role matches
	if user.Role != role {
		return "", nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// GetUserByID retrieves a user by their ID
func (s *service) GetUserByID(id string) (*entities.User, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %w", err)
	}

	user, err := s.repository.FindByID(userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *service) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}

// generateJWT creates a JWT token for a user
func (s *service) generateJWT(user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
