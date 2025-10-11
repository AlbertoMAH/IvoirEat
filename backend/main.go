package main

import (
	"fmt"
	"log"
	"net/http"
)

// corsMiddleware autorise les requêtes cross-origin
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Pour le développement. En production, soyez plus restrictif.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialiser la connexion à la base de données
	initDB()

	// Lancer la migration avant toute autre opération sur la base de données
	migrateDatabase()

	// Créer la table (ne fera rien si elle existe déjà) et insérer les données de test
	createFruitsTable()
	seedDatabase()

	// Créer un routeur et y attacher le middleware CORS
	mux := http.NewServeMux()
	mux.HandleFunc("/api/fruits", getFruitsHandler)

	handler := corsMiddleware(mux)

	fmt.Println("Server is starting on port 8080...")
	// Démarrer le serveur sur le port 8080 avec le handler CORS
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
