package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

// ValidationMiddleware gère la configuration CORS.
func ValidationMiddleware() gin.HandlerFunc {
	// Création d'une configuration CORS personnalisée avec une fonction de validation
	// pour logger les requêtes et déboguer le problème.
	config := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Log de l'origine reçue pour chaque requête
			log.Printf("CORS DEBUG: Received Origin header: [%s]", origin)

			allowedOrigins := []string{"https://frontendtest-hgiq.onrender.com", "http://localhost:3000"}

			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					log.Printf("CORS DEBUG: Origin '%s' is in the allowed list. Granting access.", origin)
					return true
				}
			}

			log.Printf("CORS DEBUG: Origin '%s' is NOT in the allowed list. Denying access.", origin)
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	return cors.New(config)
}
