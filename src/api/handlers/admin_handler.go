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

		// Register user with 'admin' role and get token
		token, _, err := service.Register(req.Username, req.Email, req.Password, entities.RoleAdmin)
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

		// Prepare response with token
		response := map[string]interface{}{
			"token": token,
		}

		return c.Status(http.StatusCreated).JSON(user.SuccessResponse{
			Status:  true,
			Message: "Admin registered successfully",
			Data:    response,
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
