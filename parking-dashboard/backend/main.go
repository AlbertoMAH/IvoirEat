package main

import (
	"gobackend/pkg/database"
	"gobackend/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connexion à la base de données
	database.ConnectDatabase()

	// Création du routeur Gin
	router := gin.Default()

	// Configuration des routes
	routes.SetupRoutes(router)

	// Démarrage du serveur
	router.Run(":8080")
}
