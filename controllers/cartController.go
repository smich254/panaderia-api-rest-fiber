package controllers

import (
	// Importar los paquetes necesarios aquí

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/database"
	// ... otros paquetes
)

func AddToCart(c *fiber.Ctx) error {
    db := database.InitDB()
    defer db.Close()

    // Estructura para representar un elemento del carrito
    var cartItem struct {
        UserID    int `json:"userID"`
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

    // Verificar si userID y productID son válidos
    var userCount, productCount int
    db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", cartItem.UserID).Scan(&userCount)
    db.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", cartItem.ProductID).Scan(&productCount)

    if userCount == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userID"})
    }

    if productCount == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid productID"})
    }

    // Finalmente, insertar en la tabla de carritos
    _, err := db.Exec("INSERT INTO carts (userID, productID, quantity) VALUES (?, ?, ?)",
        cartItem.UserID, cartItem.ProductID, cartItem.Quantity)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not insert item into cart"})
    }

    return c.JSON(fiber.Map{"message": "Product added to cart"})
}


func UpdateCartItem(c *fiber.Ctx) error {
    db := database.InitDB()
    defer db.Close()

    // Estructura para representar un elemento del carrito
    var cartItem struct {
        CartID    int `json:"cartID"`
        Quantity  int `json:"quantity"`
    }

    if err := c.BodyParser(&cartItem); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    // Verificar si la cantidad es un número positivo
    if cartItem.Quantity <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Quantity must be a positive number"})
    }

    // Verificar si cartID es válido
    var cartCount int
    db.QueryRow("SELECT COUNT(*) FROM carts WHERE id = ?", cartItem.CartID).Scan(&cartCount)

    if cartCount == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cartID"})
    }

    // Finalmente, actualizar en la tabla de carritos
    _, err := db.Exec("UPDATE carts SET quantity = ? WHERE id = ?",
        cartItem.Quantity, cartItem.CartID)
    if err != nil {
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

func Checkout(c *fiber.Ctx) error {
	db := database.InitDB()
	defer db.Close()

	// Obtener el userID desde el parámetro de la URL
	userID := c.Params("userID")

	// Verificar si el userID es válido
	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&userCount)
	if err != nil {
		log.Println("Error al verificar el userID:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	if userCount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userID"})
	}

	// Eliminar todos los elementos del carrito para este usuario
	_, err = db.Exec("DELETE FROM carts WHERE userID = ?", userID)
	if err != nil {
		log.Println("Error al vaciar el carrito:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not empty cart"})
	}

	return c.JSON(fiber.Map{"message": "Checkout successful, cart emptied"})
}