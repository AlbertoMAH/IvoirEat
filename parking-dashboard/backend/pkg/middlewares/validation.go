package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ValidationMiddleware gère la validation des entrées et la configuration CORS.
// Pour l'instant, il ne fait que configurer CORS.
func ValidationMiddleware() gin.HandlerFunc {
	// Utilise la configuration par défaut du middleware CORS,
	// qui est permissive et autorise toutes les origines.
	// C'est une bonne configuration pour le développement.
	return cors.Default()
}
