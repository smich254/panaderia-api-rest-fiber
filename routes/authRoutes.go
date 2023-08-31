package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/controllers" // Asegúrate de que esta ruta sea correcta
)

// SetupAuthRoutes configura las rutas de autenticación y registro
func SetupAuthRoutes(app *fiber.App) {
	app.Post("/api/login", controllers.Login)
	app.Post("/api/admin-login", controllers.AdminLogin)
	app.Post("/api/register", controllers.Register)  // Nueva línea
}
