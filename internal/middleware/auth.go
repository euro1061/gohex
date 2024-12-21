package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from cookie
		token := c.Cookies("token")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized - Please login first",
			})
		}

		// Store token in locals for later use
		c.Locals("token", token)

		return c.Next()
	}
}
