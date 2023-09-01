package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/controllers"
	middleware "github.com/smich254/panaderia-api-rest-fiber/middlewares"
	"github.com/smich254/panaderia-api-rest-fiber/routes"
)

func main() {
	app := fiber.New()

	// Configura las rutas públicas (ej. login, registro)
	routes.SetupAuthRoutes(app)

	// Usa el middleware de registro para todas las rutas
	app.Use(middleware.Logging())

	// Crea un nuevo grupo de rutas que será protegido por el middleware JWT
	protected := app.Group("/api", middleware.JWTMiddleware())

	// Configura las rutas protegidas (aquellas que requieren autenticación JWT)
	routes.SetupProtectedRoutes(protected)

	// Rutas para el carrito de compras (solo para usuarios autenticados)
	protected.Post("/cart/add", controllers.AddToCart)
	protected.Put("/cart/update", controllers.UpdateCartItem)
	protected.Delete("/cart/delete", controllers.DeleteFromCart)

	// Configura las rutas de autenticación
	routes.SetupAuthRoutes(app)
	
	// Configura las rutas de productos
	routes.SetupProductRoutes(app)


    // Iniciar la base de datos y crear tablas si no existen
	// Descomentar las 2 lineas de código para el primer uso
    //database.SetupDB()
	//database.SetupProductAndCartTables()

	// Escucha en el puerto 3000
	app.Listen(":3000")
}
