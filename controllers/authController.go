package controllers

import (
	"database/sql"
	"encoding/hex"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smich254/panaderia-api-rest-fiber/database"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("@123@") /// Cambia esto a tu propia llave secreta

func GenerateJWT(email string, isAdmin bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["admin"] = isAdmin
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(jwtSecret)

	return t, err
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	if !checkPassword(hashedPassword, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := GenerateJWT("user@example.com", false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

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

	err := db.QueryRow("SELECT email, password, is_admin FROM users WHERE email = ?", admin.Email).Scan(&email, &hashedPassword, &isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid admin credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	if !checkPassword(hashedPassword, admin.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid admin credentials"})
	}

	if !isAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not an admin"})
	}

	if admin.AdminCode != "someHardcodedAdminCode" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid admin code"})
	}

	token, err := GenerateJWT("admin@example.com", true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func Register(c *fiber.Ctx) error {
	db := database.InitDB()
	defer db.Close()

	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"is_admin"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	hashedPassword := hashPassword(user.Password)

	_, err := db.Exec("INSERT INTO users (email, password, is_admin) VALUES (?, ?, ?)", user.Email, hashedPassword, user.IsAdmin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not insert user"})
	}

	return c.JSON(fiber.Map{"message": "User registered"})
}
