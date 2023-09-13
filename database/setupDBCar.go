package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func SetupProductAndCartTables() {
	db, err := sql.Open("sqlite3", "database/panaderia.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crear la tabla de categorias
	createCategoryTableQuery := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nameCategory VARCHAR(15) NOT NULL
	);
	`
	_, err = db.Exec(createCategoryTableQuery)
	if err != nil {
		log.Fatal(err)
		log.Println("Error al escanear producto:", err)
	}

	// Crear la tabla de productos
	// Las llaves foraneas se deben colocar al final de la consulta
	createProductTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(15) NOT NULL,
		description VARCHAR(255),
		categoryID INTEGER NOT NULL,
		price REAL NOT NULL,
		stock INTEGER NOT NULL,
		imageURL TEXT,
		FOREIGN KEY (categoryID) REFERENCES categories (id)
	);
	`
	_, err = db.Exec(createProductTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Crear la tabla de carrito de compras
	createCartTableQuery := `
	CREATE TABLE IF NOT EXISTS carts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		userID INTEGER NOT NULL,
		productID INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		FOREIGN KEY (userID) REFERENCES users (id),
		FOREIGN KEY (productID) REFERENCES products (id)
	);
	`
	_, err = db.Exec(createCartTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	//createImageTableQuery := `
	//CREATE TABLE IF NOT EXISTS images (
	//	id INTEGER PRIMARY KEY AUTOINCREMENT,
	//	data BLOB
	//);
	//`
	//_, err = db.Exec(createImageTableQuery)
	//if err != nil {
	//	log.Fatal(err)
	//}

	log.Println("Tables products, carts, categories,  created successfully.")
}
