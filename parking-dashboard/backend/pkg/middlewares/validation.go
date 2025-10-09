package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ValidationMiddleware gère la configuration CORS.
func ValidationMiddleware() gin.HandlerFunc {
	// Création d'une configuration CORS personnalisée et plus explicite
	// pour s'assurer que tous les en-têtes nécessaires sont autorisés.
	config := cors.Config{
		AllowOrigins:     []string{"https://frontendtest-hgiq.onrender.com", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	return cors.New(config)
}
