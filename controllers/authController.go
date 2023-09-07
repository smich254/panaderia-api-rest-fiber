package controllers

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smich254/panaderia-api-rest-fiber/database"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // Obtener la llave secreta de una variable de entorno

func GenerateJWT(email string, isAdmin bool, isLoggedIn bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// La parte de "expiración" ahora se manejará aquí
	expTime := time.Now().Add(72 * time.Hour).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["admin"] = isAdmin
	claims["exp"] = expTime
	claims["isLoggedIn"] = isLoggedIn

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Error while generating JWT:", err)
		return "", err
	}

	return t, nil
}

func checkPassword(hashedPassword, password string) bool {
	hashedBytes, err := hex.DecodeString(hashedPassword)
	if err != nil {
		log.Printf("Error decoding hex: %v", err)
		return false
	}
	err = bcrypt.CompareHashAndPassword(hashedBytes, []byte(password))
	return err == nil
}

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	return hex.EncodeToString(hashed)
}

func Login(c *fiber.Ctx) error {
	db := database.InitDB()
	defer db.Close()

	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var email, hashedPassword string
	err := db.QueryRow("SELECT email, password FROM users WHERE email = ?", user.Email).Scan(&email, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with email: %s", user.Email) // Mensaje de depuración agregado
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	log.Printf("DB password: %s, Provided password: %s", hashedPassword, user.Password) // Mensaje de depuración agregado

	if !checkPassword(hashedPassword, user.Password) {
		log.Printf("Passwords do not match for email: %s", user.Email) // Mensaje de depuración agregado
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := GenerateJWT(email, false, true) // Corregido para utilizar la variable email
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}
	fmt.Println("User logged in successfully")
	return c.JSON(fiber.Map{"token": token})
}

func AdminLogin(c *fiber.Ctx) error {
	db := database.InitDB()
	defer db.Close()

	var admin struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		AdminCode string `json:"admin_code"`
	}

	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var email, hashedPassword string
	var isAdmin bool

	err := db.QueryRow("SELECT email, password, isAdmin FROM users WHERE email = ?", admin.Email).Scan(&email, &hashedPassword, &isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with email: %s", admin.Email)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid admin credentials"})
		} else {
			log.Printf("SQL Query Error: %v", err)                              // Log the SQL error
			log.Printf("Failed admin login attempt for email: %s", admin.Email) // Debug message
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	if !checkPassword(hashedPassword, admin.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid admin credentials"})
	}

	if !isAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not an admin"})
	}

	if admin.AdminCode != "@123@" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid admin code"})
	}

	token, err := GenerateJWT("admin@example.com", true, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}
	fmt.Println("User Admin logged in successfully")
	return c.JSON(fiber.Map{"token": token})
}

func Register(c *fiber.Ctx) error {
	db := database.InitDB()
	defer db.Close()

	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	hashedPassword := hashPassword(user.Password)

	_, err := db.Exec("INSERT INTO users (email, password, isAdmin) VALUES (?, ?, ?)", user.Email, hashedPassword, user.IsAdmin)
	if err != nil {
		log.Println("Error inserting new user:", err)
		log.Println("SQL Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not insert user"})
	}

	return c.JSON(fiber.Map{"message": "User registered"})
}

func Logout(c *fiber.Ctx) error {
	// Generar un nuevo token con tiempo de expiración corto (1 segundo)
	token := jwt.New(jwt.SigningMethodHS256)
	expTime := time.Now().Add(1 * time.Second).Unix() // Expira en 1 segundo

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expTime
	claims["isLoggedIn"] = false // Establecer en false

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Error while generating JWT:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not logout"})
	}

	return c.JSON(fiber.Map{"message": "Logged out", "token": t})
}
