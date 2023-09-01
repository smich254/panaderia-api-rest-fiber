package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func SetupProductsDB() {
	db, err := sql.Open("sqlite3", "database/panaderia.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crear la tabla de productos
	createTableQuery := `
	CREATE TABLE products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		description TEXT,
		price REAL,
		stock INTEGER,
		imageURL TEXT
	);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}
