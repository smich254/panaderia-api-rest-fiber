package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/controllers"
)

// SetupProductRoutes configura las rutas para los productos
func SetupProductRoutes(app *fiber.App) {
	productRouter := app.Group("/api/products")
	productRouter.Get("/", controllers.GetAllProducts)
	productRouter.Post("/", controllers.AddProduct)
	productRouter.Delete("/:id", controllers.DeleteProduct)
	productRouter.Put("/:id", controllers.UpdateProduct)
}
