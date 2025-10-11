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
		name VARCHAR(100) NOT NULL UNIQUE,
		description TEXT
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
	fruitsToSeed := map[string]string{
		"Pomme":     "Un fruit croquant et juteux, souvent rouge ou vert.",
		"Banane":    "Un fruit jaune, doux et riche en potassium.",
		"Orange":    "Un agrume rond, riche en vitamine C.",
		"Fraise":    "Un petit fruit rouge et sucré, parfait pour les desserts.",
		"Mangue":    "Un fruit tropical à la chair orange et savoureuse.",
		"Ananas":    "Un fruit tropical avec une écorce dure et une chair sucrée.",
		"Raisin":    "Petits fruits poussant en grappes, utilisés pour le vin.",
		"Cerise":    "Petit fruit rouge, souvent consommé frais ou en clafoutis.",
		"Pêche":     "Fruit à peau duveteuse et à chair juteuse.",
		"Poire":     "Fruit en forme de cloche, à la texture granuleuse.",
	}

	query := `
	INSERT INTO fruits (name, description)
	VALUES ($1, $2)
	ON CONFLICT (name)
	DO UPDATE SET description = EXCLUDED.description;
	`

	for name, description := range fruitsToSeed {
		_, err := DB.Exec(query, name, description)
		if err != nil {
			// On logue l'erreur mais on ne stoppe pas le programme
			log.Printf("Could not insert fruit %s: %v\n", name, err)
		}
	}
	fmt.Println("Database has been seeded with initial fruit data.")
}

// migrateDatabase vérifie si la colonne 'description' existe et l'ajoute si ce n'est pas le cas.
func migrateDatabase() {
	var columnExists bool
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_name = 'fruits' AND column_name = 'description'
	);`
	err := DB.QueryRow(query).Scan(&columnExists)
	if err != nil {
		log.Fatalf("Failed to check for description column: %v", err)
	}

	if !columnExists {
		fmt.Println("Column 'description' not found. Altering table...")
		alterQuery := `ALTER TABLE fruits ADD COLUMN description TEXT;`
		_, err := DB.Exec(alterQuery)
		if err != nil {
			log.Fatalf("Failed to add description column: %v", err)
		}
		fmt.Println("Table 'fruits' has been successfully altered.")
	} else {
		fmt.Println("Column 'description' already exists. No migration needed.")
	}
}
