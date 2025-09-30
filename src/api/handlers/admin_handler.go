package handlers

import (
	"belimang/src/pkg/entities"
	"belimang/src/pkg/user"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// RegisterAdmin handles admin registration
// @Summary      Register a new admin
// @Description  Create a new admin account with role 'admin'
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Param        request  body      user.RegisterRequest  true  "Registration details"
// @Success      201      {object}  user.SuccessResponse
// @Failure      400      {object}  user.ErrorResponse
// @Failure      409      {object}  user.ErrorResponse
// @Failure      500      {object}  user.ErrorResponse
// @Router       /admin/register [post]
func RegisterAdmin(service user.Service) fiber.Handler {
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

		// Register user with 'admin' role
		newAdmin, err := service.Register(req.Username, req.Email, req.Password, entities.RoleAdmin)
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
		adminResponse := user.UserResponse{
			ID:        newAdmin.ID,
			Username:  newAdmin.Username,
			Email:     newAdmin.Email,
			Role:      newAdmin.Role,
			CreatedAt: newAdmin.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		return c.Status(http.StatusCreated).JSON(user.SuccessResponse{
			Status:  true,
			Message: "Admin registered successfully",
			Data:    adminResponse,
		})
	}
}

// LoginAdmin handles admin login
// @Summary      Admin login
// @Description  Authenticate an admin and return JWT token
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Param        request  body      user.LoginRequest  true  "Login credentials"
// @Success      200      {object}  user.SuccessResponse
// @Failure      400      {object}  user.ErrorResponse
// @Failure      401      {object}  user.ErrorResponse
// @Failure      500      {object}  user.ErrorResponse
// @Router       /admin/login [post]
func LoginAdmin(service user.Service) fiber.Handler {
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

		// Authenticate admin with 'admin' role
		token, authenticatedAdmin, err := service.Login(req.Username, req.Password, entities.RoleAdmin)
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
				ID:        authenticatedAdmin.ID,
				Username:  authenticatedAdmin.Username,
				Email:     authenticatedAdmin.Email,
				Role:      authenticatedAdmin.Role,
				CreatedAt: authenticatedAdmin.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
		}

		return c.Status(http.StatusOK).JSON(user.SuccessResponse{
			Status:  true,
			Message: "Login successful",
			Data:    loginResponse,
		})
	}
}

// GetCurrentAdmin handles GET /admin/me - returns current authenticated admin info
// @Summary      Get current admin
// @Description  Get authenticated admin information
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  user.SuccessResponse
// @Failure      401      {object}  user.ErrorResponse
// @Failure      403      {object}  user.ErrorResponse
// @Failure      404      {object}  user.ErrorResponse
// @Failure      500      {object}  user.ErrorResponse
// @Router       /admin/me [get]
func GetCurrentAdmin(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context (set by JWT middleware)
		userID := c.Locals("user_id")

		// Ensure it's a string
		userUUID, ok := userID.(string)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Internal server error",
				Error:   "invalid user_id type",
			})
		}

		// Get user from service
		currentAdmin, err := service.GetUserByID(userUUID)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err == user.ErrUserNotFound {
				statusCode = http.StatusNotFound
			}

			return c.Status(statusCode).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Failed to get admin",
				Error:   err.Error(),
			})
		}

		// Prepare response
		adminResponse := user.UserResponse{
			ID:        currentAdmin.ID,
			Username:  currentAdmin.Username,
			Email:     currentAdmin.Email,
			Role:      currentAdmin.Role,
			CreatedAt: currentAdmin.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		return c.Status(http.StatusOK).JSON(user.SuccessResponse{
			Status:  true,
			Message: "Admin retrieved successfully",
			Data:    adminResponse,
		})
	}
}
