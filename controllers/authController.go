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

func GenerateJWT(name string, lastName string, userName string, email string, isAdmin bool, isLoggedIn bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// La parte de "expiración" ahora se manejará aquí
	expTime := time.Now().Add(72 * time.Hour).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["lastName"] = lastName
	claims["userName"] = userName
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
		UserName    string `json:"userName"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var email, hashedPassword, name, lastName, userName string
	err := db.QueryRow("SELECT name, lastName, userName, email, password, FROM users WHERE userName = ?", user.UserName).Scan(&name, &lastName, &userName, &email, &hashedPassword,  )
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with userName: %s", user.UserName) // Mensaje de depuración agregado
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	log.Printf("DB password: %s, Provided password: %s", hashedPassword, user.Password) // Mensaje de depuración agregado

	if !checkPassword(hashedPassword, user.Password) {
		log.Printf("Passwords do not match for userName: %s", user.UserName) // Mensaje de depuración agregado
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credSentials"})
	}

	token, err := GenerateJWT(name, lastName, userName, email, false, true) // Corregido para utilizar la variable email
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
		UserName	string `json:"userName"`
		Password  string `json:"password"`
		AdminCode string `json:"admin_code"`
	}

	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var email, hashedPassword, name, lastName, userName string
	var isAdmin bool

	
	err := db.QueryRow("SELECT name, lastName, userName, email, password, isAdmin FROM users WHERE userName = ?", admin.UserName).Scan(&name, &lastName, &userName, &email, &hashedPassword, &isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with userName: %s", admin.UserName)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid admin credentials"})
		} else {
			log.Printf("SQL Query Error: %v", err)                              // Log the SQL error
			log.Printf("Failed admin login attempt for userName: %s", admin.UserName) // Debug message
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

	token, err := GenerateJWT(name, lastName, userName, email, true, true)
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
		Name 		string `json:"name"`
		LastName 	string `json:"lastName"`
		UserName	string `json:"userName"`
		Email    	string `json:"email"`
		Password 	string `json:"password"`
		IsAdmin  	bool   `json:"isAdmin"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	hashedPassword := hashPassword(user.Password)

	_, err := db.Exec("INSERT INTO users (name, lastName, userName, email, password, isAdmin) VALUES (?, ?, ?, ?, ?, ?)", user.Name, user.LastName, user.UserName, user.Email, hashedPassword, user.IsAdmin)
	if err != nil {
		log.Println("Error inserting new user:", err)
		log.Println("SQL Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not insert user"})
	}

	return c.JSON(fiber.Map{"message": "User registered"})
}

func Logout(c *fiber.Ctx) error {
	// Obtener el token del contexto local
	user, ok := c.Locals("user").(*jwt.Token)
	if ok {
		claims := user.Claims.(jwt.MapClaims)
		claims["isLoggedIn"] = false // Establecer en false
	}

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


func isAdmin(c *fiber.Ctx) bool {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["admin"].(bool)
	return admin
}