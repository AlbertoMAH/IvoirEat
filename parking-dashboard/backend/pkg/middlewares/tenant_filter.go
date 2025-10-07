package middlewares

import "github.com/gin-gonic/gin"

// TenantFilterMiddleware filtre les requêtes pour s'assurer qu'un utilisateur n'accède qu'aux ressources de son propre tenant.
func TenantFilterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logique de filtrage par tenant_id à implémenter ici.
		// Par exemple, vérifier que le tenant_id du JWT correspond à celui de la ressource demandée.
		c.Next()
	}
}
