package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LogoutMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*jwt.Token)
		if !ok {
			return c.Next()
		}

		claims := user.Claims.(jwt.MapClaims)
		claims["isLoggedIn"] = false // Establecer en false

		return c.Next()
	}
}
