package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware valide le token JWT et injecte les informations de l'utilisateur dans le contexte.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logique de validation du token JWT à implémenter ici.
		// Pour l'instant, on laisse passer toutes les requêtes.
		c.Next()
	}
}

// RoleAuthMiddleware vérifie si l'utilisateur a le rôle requis.
func RoleAuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logique de vérification du rôle à implémenter ici.
		userRole, exists := c.Get("userRole") // Supposons que userRole est injecté par JWTMiddleware

		if !exists || userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
			c.Abort()
			return
		}
		c.Next()
	}
}
