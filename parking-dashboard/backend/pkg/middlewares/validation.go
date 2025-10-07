package middlewares

import "github.com/gin-gonic/gin"

// ValidationMiddleware gère la validation des entrées et la configuration CORS.
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Configuration CORS
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Logique de validation des entrées à implémenter ici si nécessaire.
		c.Next()
	}
}
