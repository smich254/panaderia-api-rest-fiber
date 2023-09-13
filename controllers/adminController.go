package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/database"
)

type User struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	LastName	string  `json:"lastName"`
	Email		string	`json:"email"`
	Password	string 	`json:"password"`
	IsAdmin     bool    `json:"isAdmin"`
}

// GetAllUsersByAdmin obtiene todos los usuarios
func GetAllUsersByAdmin(c *fiber.Ctx) error {
	log.Println("Fetching all users...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Password, &user.IsAdmin)
		if err != nil {
			log.Println("Error al escanear usuarios:", err)
			continue
		}
		users = append(users, user)
	}
	return c.JSON(users)
}

// AddPUserByAdmin agrega un nuevo usuario por medio de adminuser
func AddPUserByAdmin(c *fiber.Ctx) error {
	log.Println("Adding a new product...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	_, err := db.Exec("INSERT INTO users (name, lastName, email, password, isAdmin) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.LastName, user.Email, user.Password, user.IsAdmin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not insert user"})
	}

	return c.JSON(fiber.Map{"message": "User added"})
}


// DeleteUser elimina un usuario por su ID
func DeleteUserByAdmin(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted"})
}

// UpdateUserByAdmin actualiza un usuario por su ID, desde adminuser
func UpdateUserByAdmin(c *fiber.Ctx) error {
	log.Println("Updating a user...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	_, err = db.Exec("UPDATE users SET name = ?, lastName = ?, email = ?, password = ?, isAdmin = ?",
		user.Name, user.LastName, user.Email, user.Password, user.IsAdmin, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
	}

	return c.JSON(fiber.Map{"message": "User updated"})
}