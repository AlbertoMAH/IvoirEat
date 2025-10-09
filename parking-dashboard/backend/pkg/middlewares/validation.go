package middlewares

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ValidationMiddleware gère la configuration CORS de manière dynamique.
func ValidationMiddleware() gin.HandlerFunc {
	// Récupérer les origines autorisées depuis la variable d'environnement.
	// C'est une chaîne de caractères séparée par des virgules.
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	var config cors.Config

	// Si la variable d'environnement n'est pas définie, on utilise une configuration
	// permissive pour le développement local.
	if allowedOrigins == "" {
		config = cors.DefaultConfig() // DefaultConfig est déjà assez permissif
		config.AllowAllOrigins = true // On le rend encore plus explicite
		config.AllowCredentials = true
		// On s'assure que les en-têtes nécessaires à l'authentification JWT sont présents
		config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	} else {
		// Si la variable est définie, on l'utilise pour configurer CORS de manière stricte.
		config = cors.Config{
			// On divise la chaîne de caractères en une slice de strings
			AllowOrigins:     strings.Split(allowedOrigins, ","),
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "Cache-Control", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}
	}

	// Retourne le nouveau middleware CORS avec la configuration choisie.
	return cors.New(config)
}
