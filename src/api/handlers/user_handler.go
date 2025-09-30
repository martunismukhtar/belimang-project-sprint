package handlers

import (
	"belimang/src/pkg/entities"
	"belimang/src/pkg/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// RegisterUser handles user registration
// @Summary      Register a new user
// @Description  Create a new user account with role 'user'
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      user.RegisterRequest  true  "Registration details"
// @Success      201      {object}  user.SuccessResponse
// @Failure      400      {object}  user.ErrorResponse
// @Failure      409      {object}  user.ErrorResponse
// @Failure      500      {object}  user.ErrorResponse
// @Router       /users/register [post]
func RegisterUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req user.RegisterRequest

		// Parse request body
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Invalid request body",
				Error:   err.Error(),
			})
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Validation failed",
				Error:   err.Error(),
			})
		}

		// Register user with 'user' role
		newUser, err := service.Register(req.Username, req.Email, req.Password, entities.RoleUser)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err == user.ErrUsernameExists || err == user.ErrEmailExists {
				statusCode = http.StatusConflict
			} else if err == user.ErrInvalidEmail || err == user.ErrInvalidUsername || err == user.ErrInvalidPassword {
				statusCode = http.StatusBadRequest
			}

			return c.Status(statusCode).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Registration failed",
				Error:   err.Error(),
			})
		}

		// Prepare response
		userResponse := user.UserResponse{
			ID:        newUser.ID,
			Username:  newUser.Username,
			Email:     newUser.Email,
			Role:      newUser.Role,
			CreatedAt: newUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		return c.Status(http.StatusCreated).JSON(user.SuccessResponse{
			Status:  true,
			Message: "User registered successfully",
			Data:    userResponse,
		})
	}
}

// LoginUser handles user login
// @Summary      User login
// @Description  Authenticate a user and return JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      user.LoginRequest  true  "Login credentials"
// @Success      200      {object}  user.SuccessResponse
// @Failure      400      {object}  user.ErrorResponse
// @Failure      401      {object}  user.ErrorResponse
// @Failure      500      {object}  user.ErrorResponse
// @Router       /users/login [post]
func LoginUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req user.LoginRequest

		// Parse request body
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Invalid request body",
				Error:   err.Error(),
			})
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Validation failed",
				Error:   err.Error(),
			})
		}

		// Authenticate user with 'user' role
		token, authenticatedUser, err := service.Login(req.Username, req.Password, entities.RoleUser)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err == user.ErrInvalidCredentials {
				statusCode = http.StatusUnauthorized
			}

			return c.Status(statusCode).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Login failed",
				Error:   err.Error(),
			})
		}

		// Prepare response
		loginResponse := user.LoginResponse{
			Token: token,
			User: user.UserResponse{
				ID:        authenticatedUser.ID,
				Username:  authenticatedUser.Username,
				Email:     authenticatedUser.Email,
				Role:      authenticatedUser.Role,
				CreatedAt: authenticatedUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
		}

		return c.Status(http.StatusOK).JSON(user.SuccessResponse{
			Status:  true,
			Message: "Login successful",
			Data:    loginResponse,
		})
	}
}