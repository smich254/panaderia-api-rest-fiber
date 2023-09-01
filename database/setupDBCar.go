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

	// Crear la tabla de productos
	createProductTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		stock INTEGER NOT NULL,
		imageURL TEXT
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

	log.Println("Tables products and carts created successfully.")
}
