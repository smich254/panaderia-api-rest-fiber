package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/middlewares"
	"github.com/smich254/panaderia-api-rest-fiber/routes"
)

func main() {
	app := fiber.New()

	// Configura las rutas públicas (ej. login, registro)
	routes.SetupAuthRoutes(app)

	// Usa el middleware de registro para todas las rutas
	app.Use(middlewares.Logging())

	// Usa el middleware de logout para todas las rutas
	app.Use(middlewares.LogoutMiddleware())

	// Crea un nuevo grupo de rutas que será protegido por el middleware JWT

	// Configura las rutas protegidas (aquellas que requieren autenticación JWT)
	routes.SetupProtectedRoutes(app)

    // Configura las rutas del carrito de compras
    routes.SetupCartRoutes(app)

	// Configura las rutas de productos
	//routes.SetupProductRoutes(app)

	// Iniciar la base de datos y crear tablas si no existen
	// Descomentar las 2 lineas de código para el primer uso
	// Nota: Actualizar las herramientas de go desde VS Code antes
	// De descomentar
	//database.SetupDB()
	//database.SetupProductAndCartTables()

	// Escucha en el puerto 3000

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		  "error": "Error 404: La página que estás buscando no se pudo encontrar. Si eres un desarrollador, consulta la documentación de la API para rutas válidas.",
		})
	  })
	  
	app.Listen(":3000")
}
