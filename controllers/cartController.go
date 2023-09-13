package controllers

import (
	// Importar los paquetes necesarios aquí

	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/smich254/panaderia-api-rest-fiber/database"
	// ... otros paquetes
)

func AddToCart(c *fiber.Ctx) error {
    db := database.InitDB()
    defer db.Close()

    var cartItem struct {
        ProductID int `json:"productID"`
        Quantity  int `json:"quantity"`
    }

    if err := c.BodyParser(&cartItem); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    // Verificar si la cantidad es un número positivo
    if cartItem.Quantity <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Quantity must be a positive number"})
    }

    tokenStr := c.Get("Authorization")[7:] // Suponiendo que el encabezado es "Bearer <token>"

    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte("tu_clave_secreta"), nil // Reemplaza "tu_clave_secreta" con tu clave secreta real
    })

    if err != nil {
        log.Println("Error al decodificar el token:", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        log.Println("Error al obtener las claims del token")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    userID := int(claims["userID"].(float64)) // Suponiendo que el userID se almacena como un número flotante en el token

    // Verificar si productID es válido
    var productCount int
    err = db.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", cartItem.ProductID).Scan(&productCount)
    if err != nil {
        log.Println("Error al contar los productos:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
    }

    if productCount == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid productID"})
    }

    _, err = db.Exec("INSERT INTO carts (userID, productID, quantity) VALUES (?, ?, ?)",
        userID, cartItem.ProductID, cartItem.Quantity)
    if err != nil {
        log.Println("Error al insertar el elemento en el carrito:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not insert item into cart"})
    }

    return c.JSON(fiber.Map{"message": "Product added to cart"})
}



func UpdateCartItem(c *fiber.Ctx) error {
    db := database.InitDB()
    defer db.Close()

    // Estructura para representar un elemento del carrito
    var cartItem struct {
        Quantity  int `json:"quantity"`
    }

    if err := c.BodyParser(&cartItem); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    // Verificar si la cantidad es un número positivo
    if cartItem.Quantity <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Quantity must be a positive number"})
    }

    tokenStr := c.Get("Authorization")[7:] // Suponiendo que el encabezado es "Bearer <token>"

    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte("tu_clave_secreta"), nil // Reemplaza "tu_clave_secreta" con tu clave secreta real
    })

    if err != nil {
        log.Println("Error al decodificar el token:", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        log.Println("Error al obtener las claims del token")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    userID := int(claims["userID"].(float64)) // Suponiendo que el userID se almacena como un número flotante en el token

    // Obtener el cartID basado en el userID
    var cartID int
    err = db.QueryRow("SELECT id FROM carts WHERE userID = ?", userID).Scan(&cartID)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userID or cart not found"})
        }
        log.Println("Error al consultar el cartID:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
    }

    // Verificar si cartID es válido
    var cartCount int
    err = db.QueryRow("SELECT COUNT(*) FROM carts WHERE id = ?", cartID).Scan(&cartCount)
    if err != nil {
        log.Println("Error al contar los elementos del carrito:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
    }

    if cartCount == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cartID"})
    }

    _, err = db.Exec("UPDATE carts SET quantity = ? WHERE id = ?", cartItem.Quantity, cartID)
    if err != nil {
        log.Println("Error al actualizar el elemento del carrito:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update cart item"})
    }

    return c.JSON(fiber.Map{"message": "Cart item updated"})
}



func DeleteFromCart(c *fiber.Ctx) error {
    db := database.InitDB()
    defer db.Close()

    cartID := c.Params("cartID")

    // Verificar si cartID es válido
    var cartCount int
    db.QueryRow("SELECT COUNT(*) FROM carts WHERE id = ?", cartID).Scan(&cartCount)

    if cartCount == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cartID"})
    }

    // Finalmente, eliminar de la tabla de carritos
    _, err := db.Exec("DELETE FROM carts WHERE id = ?", cartID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete cart item"})
    }

    return c.JSON(fiber.Map{"message": "Product removed from cart"})
}

