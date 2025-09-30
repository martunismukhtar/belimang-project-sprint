package user

import (
	"belimang/src/pkg/entities"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrUsernameExists        = errors.New("username already exists")
	ErrEmailExists           = errors.New("email already exists for this role")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidUsername       = errors.New("username must be 5-30 characters and contain only alphanumeric characters and underscores")
	ErrInvalidPassword       = errors.New("password must be 5-30 characters")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidRole           = errors.New("invalid role")
	emailRegex               = regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	usernameRegex            = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// Service interface defines business logic operations for users
type Service interface {
	Register(username, email, password, role string) (string, *entities.User, error)
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
		panic("JWT_SECRET environment variable is required for security. Please set it in your environment or .env file")
	}
	return &service{
		repository: repo,
		jwtSecret:  jwtSecret,
	}
}

// Register creates a new user account with validation and returns a JWT token
func (s *service) Register(username, email, password, role string) (string, *entities.User, error) {
	// Validate role
	if role != entities.RoleUser && role != entities.RoleAdmin {
		return "", nil, ErrInvalidRole
	}

	// Sanitize and validate username
	username = strings.TrimSpace(username)
	usernameLen := utf8.RuneCountInString(username)
	if usernameLen < 5 || usernameLen > 30 {
		return "", nil, ErrInvalidUsername
	}
	// Only allow alphanumeric characters and underscore
	if !usernameRegex.MatchString(username) {
		return "", nil, ErrInvalidUsername
	}

	// Sanitize and validate email
	email = strings.TrimSpace(strings.ToLower(email))
	if !emailRegex.MatchString(email) {
		return "", nil, ErrInvalidEmail
	}

	// Validate password length (plain text before hashing)
	if len(password) < 5 || len(password) > 30 {
		return "", nil, ErrInvalidPassword
	}

	// Hash password before attempting to create user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user - let database unique constraints handle duplicates
	user := &entities.User{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
		Role:     role,
		Password: string(hashedPassword),
	}

	if err := s.repository.Create(user); err != nil {
		// Check if it's a unique constraint violation
		errMsg := err.Error()
		if strings.Contains(errMsg, "duplicate key") || strings.Contains(errMsg, "UNIQUE constraint") || strings.Contains(errMsg, "unique_violation") {
			if strings.Contains(errMsg, "username") || strings.Contains(errMsg, "idx_users_username_unique") {
				return "", nil, ErrUsernameExists
			}
			if strings.Contains(errMsg, "email") || strings.Contains(errMsg, "idx_users_email_role_unique") {
				return "", nil, ErrEmailExists
			}
		}
		return "", nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token for the newly registered user
	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
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
