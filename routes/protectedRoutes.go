package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/controllers"
	"github.com/smich254/panaderia-api-rest-fiber/middlewares"
)

// SetupProtectedRoutes configura las rutas que requieren autenticación JWT
func SetupProtectedRoutes(app *fiber.App) {

	// Crear un grupo de rutas para el usuario
	userGroupProtected := app.Group("/api/user", middlewares.JWTMiddleware())
	
	// Ruta para obtener el perfil del usuario
	userGroupProtected.Get("/profile", controllers.GetUserProfile)

	// Ruta para actualizar el perfil del usuario
	userGroupProtected.Put("/profile", controllers.UpdateUserProfile)

	// Crear un grupo de rutas para el administrador
	adminGroupProtected := app.Group("/api/admin", middlewares.JWTMiddleware())

	// Ruta para listar todos los usuarios por medio de adminuser
	adminGroupProtected.Get("/users", controllers.GetAllUsersByAdmin)

	// Ruta para agregar un usuario por medio de AdminUser
	adminGroupProtected.Get("/adduser", controllers.AddPUserByAdmin)

	// Ruta para actualizar un usuario por medio de adminuser
	adminGroupProtected.Put("/users/:id", controllers.UpdateUserByAdmin)

	// Ruta para eliminar un usuario por medio de adminuser
	adminGroupProtected.Delete("/users/:id", controllers.DeleteUserByAdmin)

	// Ruta para obtener un producto
	adminGroupProtected.Get("/products", controllers.GetAllProducts)

	// Ruta para agregar un producto
	adminGroupProtected.Post("/products", controllers.AddProduct)

	// Ruta para actualizar un producto
	adminGroupProtected.Put("/products/:id", controllers.UpdateProduct)

	// Ruta para eliminar un producto
	adminGroupProtected.Delete("/products/:id", controllers.DeleteProduct)

	// Ruta para obtener todas las categorias
	adminGroupProtected.Get("/categories", controllers.GetAllCategories)

	// Ruta para agregar una categoria
	adminGroupProtected.Post("/categories", controllers.AddCategory)

	// Ruta para actualizar una categoria
	adminGroupProtected.Put("/categories/:id", controllers.UpdateCategory)

	// Ruta para eliminar una categoria
	adminGroupProtected.Delete("/categories/:id", controllers.DeleteCategory)
}
