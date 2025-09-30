package middleware

import (
	"belimang/src/pkg/entities"
	"belimang/src/pkg/user"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTAuth is a middleware that validates JWT tokens
func JWTAuth(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Unauthorized",
				Error:   "authorization header required",
			})
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Unauthorized",
				Error:   "invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// Validate token
		claims, err := userService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Unauthorized",
				Error:   "invalid or expired token",
			})
		}

		// Extract user info from claims
		userID, ok := (*claims)["user_id"].(string)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Unauthorized",
				Error:   "invalid token claims",
			})
		}

		// Store user info in context (as string, will be parsed in handlers)
		c.Locals("user_id", userID)
		c.Locals("username", (*claims)["username"])
		c.Locals("email", (*claims)["email"])
		c.Locals("role", (*claims)["role"])

		return c.Next()
	}
}

// IsUser middleware ensures the authenticated user has 'user' role
func IsUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(http.StatusUnauthorized).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Unauthorized",
				Error:   "role not found in context",
			})
		}

		if role != entities.RoleUser {
			return c.Status(http.StatusForbidden).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Forbidden",
				Error:   "insufficient permissions",
			})
		}

		return c.Next()
	}
}

// IsAdmin middleware ensures the authenticated user has 'admin' role
func IsAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(http.StatusUnauthorized).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Unauthorized",
				Error:   "role not found in context",
			})
		}

		if role != entities.RoleAdmin {
			return c.Status(http.StatusForbidden).JSON(user.ErrorResponse{
				Status:  false,
				Message: "Forbidden",
				Error:   "insufficient permissions",
			})
		}

		return c.Next()
	}
}
