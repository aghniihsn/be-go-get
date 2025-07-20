package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Middlewares validates JWT token and checks user role
func Middlewares(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authToken := c.Get("Authorization")
		if authToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization Header",
			})
		}

		dataDecode, err := Decoder(authToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid / Expired Authorization Token",
			})
		}

		// Check if user role is allowed (if requiredRoles specified)
		if len(requiredRoles) > 0 {
			roleAllowed := false
			for _, role := range requiredRoles {
				if dataDecode.Role == role {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "You don't have permission to access this resource",
				})
			}
		}

		// Store user data in context for use in controllers
		c.Locals("user", dataDecode)
		return c.Next()
	}
}

// AuthRequired middleware for routes that require authentication
func AuthRequired() fiber.Handler {
	return Middlewares()
}

// AdminOnly middleware for admin-only routes
func AdminOnly() fiber.Handler {
	return Middlewares("admin")
}

// UserOrAdmin middleware for routes accessible by users and admins
func UserOrAdmin() fiber.Handler {
	return Middlewares("user", "admin")
}
