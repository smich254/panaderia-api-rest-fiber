package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDB() {
	// Crear la carpeta si no existe
	if _, err := os.Stat("database"); os.IsNotExist(err) {
		err := os.Mkdir("database", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Resto del código
	db, err := sql.Open("sqlite3", "database/panaderia.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableQuery := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(20),
		lastName VARCHAR(20),
		email TEXT UNIQUE,
		password VARCHAR(20),
		isAdmin BOOLEAN
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

}
