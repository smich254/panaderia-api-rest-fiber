package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/controllers"
)

// SetupProtectedRoutes configura las rutas que requieren autenticación JWT
func SetupProtectedRoutes(app fiber.Router) {
	// Crear un grupo de rutas para el usuario
	userGroup := app.Group("/user")
	
	// Ruta para obtener el perfil del usuario
	userGroup.Get("/profile", controllers.GetUserProfile)

	// Ruta para actualizar el perfil del usuario
	userGroup.Put("/profile", controllers.UpdateUserProfile)

	// Crear un grupo de rutas para el administrador
	adminGroup := app.Group("/admin")

	// Ruta para listar todos los usuarios
	adminGroup.Get("/users", controllers.ListUsers)

	// Ruta para actualizar un usuario
	adminGroup.Put("/users/:id", controllers.UpdateUser)

	// Ruta para eliminar un usuario
	adminGroup.Delete("/users/:id", controllers.DeleteUser)

	// Ruta para agregar un producto
	adminGroup.Post("/products", controllers.AddProduct)

	// Ruta para actualizar un producto
	adminGroup.Put("/products/:id", controllers.UpdateProduct)

	// Ruta para eliminar un producto
	adminGroup.Delete("/products/:id", controllers.DeleteProduct)
}
