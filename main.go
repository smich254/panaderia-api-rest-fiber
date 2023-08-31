package main

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/smich254/panaderia-api-rest-fiber/routes"
)

// Deberías guardar esto en una variable de entorno, no en el código fuente.
var jwtSecret = []byte("@123@")

func main() {
	app := fiber.New()

	// Configura las rutas públicas (ej. login, registro)
	routes.SetupAuthRoutes(app)

	// Configura el middleware JWT
	app.Use(jwtware.New(jwtware.Config{
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			// Verifica siempre el método de firma
			if t.Method.Alg() != jwtware.HS256 {
				return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
			}
			return jwtSecret, nil
		},
	}))

	// Configura las rutas protegidas (aquellas que requieren autenticación JWT)
	// routes.SetupProtectedRoutes(app)
	// Por ejemplo:
	// app.Get("/api/user", func(c *fiber.Ctx) error {
	//    return c.SendString("Hello, authenticated user!")
	// })

	// Escucha en el puerto 3000
	app.Listen(":3000")
}
