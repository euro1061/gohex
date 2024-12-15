// internal/adapters/primary/http/middleware/auth_middleware.go
package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	// ถ้ามี dependencies เพิ่มเติมสามารถใส่ตรงนี้
	// เช่น userService, jwtSecret เป็นต้น
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) VerifyToken(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	fmt.Println(cookie)
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// ถ้าต้องการเก็บข้อมูล user ไว้ใช้ใน handler ต่อ
	c.Locals("user", claims)

	return c.Next()
}
