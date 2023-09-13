package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/smich254/panaderia-api-rest-fiber/database"
	"github.com/smich254/panaderia-api-rest-fiber/middlewares"
	"github.com/smich254/panaderia-api-rest-fiber/routes"
)

func main() {
	app := fiber.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Usa el middleware de registro para todas las rutas
	app.Use(middlewares.Logging())

	// Configura las rutas públicas (ej. login, registro)
	routes.SetupAuthRoutes(app)

	// Configura las rutas protegidas (aquellas que requieren autenticación JWT)
	routes.SetupProtectedRoutes(app)

    // Configura las rutas del carrito de compras
    routes.SetupCartRoutes(app)

	// Iniciar la base de datos y crear tablas si no existen
	// Descomentar las 2 lineas de código para el primer uso
	// Nota: Actualizar las herramientas de go desde VS Code antes
	// De descomentar
	database.SetupDB()
	
	database.SetupProductAndCartTables()

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		  "error": "Error 404: La página que estás buscando no se pudo encontrar. Si eres un desarrollador, consulta la documentación de la API para rutas válidas.",
		})
	  })

	// Escucha en el puerto 3000
	app.Listen(":3000")
}
