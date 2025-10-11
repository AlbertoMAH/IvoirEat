package main

import (
	"encoding/json"
	"net/http"
)

// getFruitsHandler gère la recherche de fruits.
// Il lit un paramètre de requête 'search' et filtre les fruits par nom.
func getFruitsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le terme de recherche depuis les paramètres de la requête URL
	searchTerm := r.URL.Query().Get("search")

	var fruits []Fruit
	var err error

	// Préparer la requête SQL. Utiliser LIKE pour une recherche insensible à la casse.
	query := `SELECT id, name FROM fruits WHERE name ILIKE $1`
	rows, err := DB.Query(query, "%"+searchTerm+"%")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Itérer sur les résultats et les ajouter à notre slice de fruits
	for rows.Next() {
		var fruit Fruit
		if err := rows.Scan(&fruit.ID, &fruit.Name); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		fruits = append(fruits, fruit)
	}

	// Définir le header pour indiquer que la réponse est du JSON
	w.Header().Set("Content-Type", "application/json")
	// Encoder la slice de fruits en JSON et l'écrire dans la réponse
	json.NewEncoder(w).Encode(fruits)
}
