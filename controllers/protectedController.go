package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/smich254/panaderia-api-rest-fiber/database"
	"github.com/smich254/panaderia-api-rest-fiber/models"
	// Importa cualquier otro paquete necesario
)

// GetUserProfile obtiene el perfil del usuario autenticado
func GetUserProfile(c *fiber.Ctx) error {
	log.Println("Fetching user profile...")

	// Middleware JWT ya debería haber almacenado el token en c.Locals("user")
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		log.Println("Unauthorized due to missing token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not authorized"})
	}

	// Obtener el ID del usuario desde las reclamaciones del token JWT
	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims["userId"].(int)
	if !ok {
		log.Println("Unauthorized due to missing userId in claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not authorized"})
	}

	db := database.InitDB()
	defer db.Close()

	row := db.QueryRow("SELECT * FROM users WHERE id = ?", userId)
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.LastName, &user.Password, &user.IsAdmin)
	if err != nil {
		log.Println("Error al escanear el perfil del usuario:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// Omitir el campo de contraseña al devolver el perfil del usuario
	user.Password = ""

	return c.JSON(user)
}

// UpdateUserProfile actualiza el perfil del usuario autenticado
func UpdateUserProfile(c *fiber.Ctx) error {
	log.Println("Updating user profile...")

	// Middleware JWT ya debería haber almacenado el token en c.Locals("user")
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		log.Println("Unauthorized due to missing token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not authorized"})
	}

	// Obtener el ID del usuario desde las reclamaciones del token JWT
	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims["userId"].(int)
	if !ok {
		log.Println("Unauthorized due to missing userId in claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not authorized"})
	}

	db := database.InitDB()
	defer db.Close()

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Verificar si el usuario intenta actualizar su propio perfil
	if user.ID != userId {
		log.Println("Unauthorized update of another user's profile")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not authorized"})
	}

	// Realizar la actualización del perfil del usuario
	_, err := db.Exec("UPDATE users SET name = ?, lastName = ?, userName = ?, email = ?, password = ? WHERE id = ?",
		user.Name, user.LastName, user.UserName, user.Email, user.Password, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user profile"})
	}

	return c.JSON(fiber.Map{"message": "User profile updated"})
}

