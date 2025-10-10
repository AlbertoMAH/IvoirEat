package main

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 1. Définir les routes de l'API en premier.
	// Le routeur de Gin donnera la priorité à ces routes spécifiques.
	api := r.Group("/api")
	{
		api.GET("/message", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Bonjour depuis le backend Go !",
			})
		})
	}

	// 2. Ensuite, utiliser le middleware pour servir les fichiers statiques.
	// Ce middleware ne sera appelé que pour les requêtes qui ne correspondent pas aux routes API.
	// Le second paramètre `true` active le mode "HTML5", qui redirige les 404 vers /index.html,
	// ce qui est parfait pour une Single Page Application.
	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	r.Run(":8080")
}
