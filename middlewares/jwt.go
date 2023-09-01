package middleware

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret debe ser la misma llave utilizada al firmar los tokens.
var jwtSecret = []byte("@123@")

// JWTMiddleware configura el middleware JWT.
func JWTMiddleware() func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			// Verificar siempre el método de firma
			if t.Method.Alg() != jwtware.HS256 {
				return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
			}
			return jwtSecret, nil
		},
	})
}
