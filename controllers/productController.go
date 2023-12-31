package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smich254/panaderia-api-rest-fiber/database"
	"github.com/smich254/panaderia-api-rest-fiber/models"
)

func GetAllCategories(c *fiber.Ctx) error {
	log.Println("Fetching all categories...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.NameCategory)
		if err != nil {
			log.Println("Error al escanear categoria:", err)
			continue
		}
		categories = append(categories, category)
	}

	return c.JSON(categories)
}

// AddCategory agrega una nueva categoria
func AddCategory(c *fiber.Ctx) error {
	log.Println("Adding a new category...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	_, err := db.Exec("INSERT INTO categories (nameCategory) VALUES (?)",
		category.NameCategory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not insert category"})
	}

	return c.JSON(fiber.Map{"message": "Category added"})
}

// DeleteProduct elimina un producto por su ID
func DeleteCategory(c *fiber.Ctx) error {
	log.Println("Deleting a category...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	id := c.Params("id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	_, err = db.Exec("DELETE FROM categories WHERE id = ?", categoryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete category"})
	}

	return c.JSON(fiber.Map{"message": "Category deleted"})
}

// UpdateCategory actualiza una categoria por su ID
func UpdateCategory(c *fiber.Ctx) error {
	log.Println("Updating a category...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	id := c.Params("id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	_, err = db.Exec("UPDATE categories SET nameCategory = ? WHERE id = ?",
		category.NameCategory, categoryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update category"})
	}

	return c.JSON(fiber.Map{"message": "Category updated"})
}

// GetAllProducts obtiene todos los productos
func GetAllProducts(c *fiber.Ctx) error {
	log.Println("Fetching all products...")
	db := database.InitDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.CategoryID, &product.Price, &product.Stock, &product.ImageURL)
		if err != nil {
			log.Println("Error al escanear producto:", err)
			continue
		}
		products = append(products, product)
	}

	return c.JSON(products)
}

// AddProduct agrega un nuevo producto
func AddProduct(c *fiber.Ctx) error {
	log.Println("Adding a new product...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	_, err := db.Exec("INSERT INTO products (name, description, categoryID, price, stock, imageURL) VALUES (?, ?, ?, ?, ?, ?)",
		product.Name, product.Description, product.CategoryID, product.Price, product.Stock, product.ImageURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not insert product"})
	}

	return c.JSON(fiber.Map{"message": "Product added"})
}

// DeleteProduct elimina un producto por su ID
func DeleteProduct(c *fiber.Ctx) error {
	log.Println("Deleting a product...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	id := c.Params("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	_, err = db.Exec("DELETE FROM products WHERE id = ?", productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete product"})
	}

	return c.JSON(fiber.Map{"message": "Product deleted"})
}

// UpdateProduct actualiza un producto por su ID
func UpdateProduct(c *fiber.Ctx) error {
	log.Println("Updating a product...")
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	db := database.InitDB()
	defer db.Close()

	id := c.Params("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	_, err = db.Exec("UPDATE products SET name = ?, description = ?, categoryID = ?, price = ?, stock = ?, imageURL = ? WHERE id = ?",
		product.Name, product.Description, product.Price, product.Stock, product.ImageURL, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update product"})
	}

	return c.JSON(fiber.Map{"message": "Product updated"})
}
