package main

import (
	"fmt"
	"log"
)

// createFruitsTable crée la table 'fruits' si elle n'existe pas.
func createFruitsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS fruits (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL UNIQUE
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create fruits table: %v", err)
	}
	fmt.Println("Table 'fruits' is ready.")
}

// seedDatabase insère des données de test dans la table 'fruits'.
// Elle ignore les erreurs si les fruits existent déjà (contrainte UNIQUE sur le nom).
func seedDatabase() {
	fruitsToSeed := []string{"Pomme", "Banane", "Orange", "Fraise", "Mangue", "Ananas"}

	query := `INSERT INTO fruits (name) VALUES ($1) ON CONFLICT (name) DO NOTHING;`

	for _, fruitName := range fruitsToSeed {
		_, err := DB.Exec(query, fruitName)
		if err != nil {
			// On logue l'erreur mais on ne stoppe pas le programme
			log.Printf("Could not insert fruit %s: %v\n", fruitName, err)
		}
	}
	fmt.Println("Database has been seeded with initial fruit data.")
}
