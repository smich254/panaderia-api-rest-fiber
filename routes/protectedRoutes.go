package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/controllers"
	"github.com/smich254/panaderia-api-rest-fiber/middlewares"
)

// SetupProtectedRoutes configura las rutas que requieren autenticación JWT
func SetupProtectedRoutes(app fiber.Router) {
	// Crear un grupo de rutas para el usuario
	userGroupProtected := app.Group("/api/user", middlewares.JWTMiddleware())
	
	// Ruta para obtener el perfil del usuario
	userGroupProtected.Get("/profile", controllers.GetUserProfile)

	// Ruta para actualizar el perfil del usuario
	userGroupProtected.Put("/profile", controllers.UpdateUserProfile)

	// Crear un grupo de rutas para el administrador
	adminGroupProtected := app.Group("/api/admin", middlewares.JWTMiddleware())

	// Ruta para listar todos los usuarios
	adminGroupProtected.Get("/users", controllers.ListUsers)

	// Ruta para actualizar un usuario
	adminGroupProtected.Put("/users/:id", controllers.UpdateUser)

	// Ruta para eliminar un usuario
	adminGroupProtected.Delete("/users/:id", controllers.DeleteUser)

	// Ruta para obtener un producto
	adminGroupProtected.Get("/products", controllers.AddProduct)

	// Ruta para agregar un producto
	adminGroupProtected.Post("/products", controllers.AddProduct)

	// Ruta para actualizar un producto
	adminGroupProtected.Put("/products/:id", controllers.UpdateProduct)

	// Ruta para eliminar un producto
	adminGroupProtected.Delete("/products/:id", controllers.DeleteProduct)
}
