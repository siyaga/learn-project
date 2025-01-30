package middleware

import (
	"learn_project/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the JWT token from the Authorization header
		tokenString := c.Get("Authorization")

		// If the Authorization header is missing or empty, return Unauthorized
		if tokenString == "" {
			return utils.ResponseError(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		// Remove "Bearer " prefix from the token string (if present)
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[len("Bearer "):]
		}

		// Validate the token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return utils.ResponseError(c, fiber.StatusUnauthorized, "Invalid token", nil)
		}

		// Add the user email to the context for use in controllers
		c.Locals("email", claims.Email)

		// Continue to the next handler
		return c.Next()
	}
}
