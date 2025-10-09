package main

import (
	"gobackend/pkg/database"
	"gobackend/pkg/models"
	"gobackend/pkg/routes"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// StripAppPrefix est un middleware pour supprimer le préfixe /app ajouté par Render.
func StripAppPrefix() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Vérifie si le chemin de la requête commence par /app
		if strings.HasPrefix(c.Request.URL.Path, "/app") {
			// Réécrit le chemin de l'URL sans le préfixe
			c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/app")
		}
		// Passe à la prochaine fonction middleware/handler
		c.Next()
	}
}

func main() {
	// Charger les variables d'environnement depuis le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Connexion à la base de données
	database.ConnectDatabase()

	// Créer l'utilisateur super_admin si nécessaire
	createSuperAdmin()

	// Création du routeur Gin
	router := gin.Default()

	// Gérer automatiquement les slashs finaux
	router.RedirectTrailingSlash = true

	// Appliquer le middleware pour le préfixe de Render
	router.Use(StripAppPrefix())

	// Configuration des routes
	routes.SetupRoutes(router)

	// Démarrage du serveur
	router.Run(":8080")
}

func createSuperAdmin() {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		log.Println("ADMIN_EMAIL or ADMIN_PASSWORD not set, skipping super admin creation")
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", adminEmail).First(&user).Error; err == nil {
		// L'utilisateur existe déjà
		log.Println("Super admin user already exists.")
		return
	}

	// Hacher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Créer le nouvel utilisateur
	superAdmin := models.User{
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     "super_admin",
	}

	if err := database.DB.Create(&superAdmin).Error; err != nil {
		log.Fatalf("Failed to create super admin: %v", err)
	}

	// log.Println("Super admin user created successfully.")
}
