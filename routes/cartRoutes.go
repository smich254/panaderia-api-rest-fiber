package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/controllers"
	"github.com/smich254/panaderia-api-rest-fiber/middlewares"
)

// SetupCartRoutes configura las rutas para el carrito de compras
func SetupCartRoutes(app *fiber.App) {
    // Crea un nuevo grupo de rutas que será protegido por el middleware JWT
    protected := app.Group("/api", middlewares.JWTMiddleware())

    // Rutas para el carrito de compras (solo para usuarios autenticados)
    protected.Post("/cart/add", controllers.AddToCart)
    protected.Put("/cart/update", controllers.UpdateCartItem)
    protected.Delete("/cart/delete", controllers.DeleteFromCart)
}
