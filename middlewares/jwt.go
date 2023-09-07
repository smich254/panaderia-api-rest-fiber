package middlewares

import (
	"fmt"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret debe ser la misma llave utilizada al firmar los tokens.
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// JWTMiddleware configura el middleware JWT.
func JWTMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Primero, ejecutar el middleware JWT original
		if err := jwtware.New(jwtware.Config{
			KeyFunc: func(t *jwt.Token) (interface{}, error) {
				// Verificar siempre el método de firma
				if t.Method.Alg() != jwtware.HS256 {
					return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
				}
				return jwtSecret, nil
			},
		})(c); err != nil {
			return err
		}

		// Luego, verificar el campo "isLoggedIn"
		if token, ok := c.Locals("user").(*jwt.Token); ok {
			claims := token.Claims.(jwt.MapClaims)
			fmt.Println(claims)
			if isLoggedIn, ok := claims["isLoggedIn"].(bool); !ok || !isLoggedIn {
				fmt.Println("Unauthorized due to isLoggedIn") 
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not logged in"})
			}
		} else {
			// Manejar el caso en que el token no existe o es nil
			fmt.Println("Unauthorized due to missing token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not authorized"})
		}

		return c.Next()
	}
}
