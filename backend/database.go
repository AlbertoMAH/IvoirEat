package main

import (
	"database/sql"
	"fmt"
	"log"
	"os" // Importer le package os

	_ "github.com/lib/pq" // Le driver PostgreSQL
)

// DB est une variable globale pour la connexion à la base de données
var DB *sql.DB

// initDB initialise la connexion à la base de données PostgreSQL.
func initDB() {
	// Lire la chaîne de connexion depuis une variable d'environnement
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fournir une valeur par défaut pour le développement local si la variable n'est pas définie
		log.Println("DATABASE_URL not set, using default connection string for local development.")
		connStr = "user=postgres password=mysecretpassword dbname=postgres sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Vérifier que la connexion est bien établie
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
}
