package main

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// Fruit définit la structure pour nos données de fruits
type Fruit struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// fruits est notre "base de données" en mémoire.
var fruits = []Fruit{
	{Name: "Pomme", Description: "Un fruit croquant et sucré, souvent rouge ou vert."},
	{Name: "Banane", Description: "Un fruit jaune, doux et riche en potassium."},
	{Name: "Orange", Description: "Un agrume juteux, source de vitamine C."},
	{Name: "Fraise", Description: "Un petit fruit rouge et sucré, parfait pour les desserts."},
	{Name: "Raisin", Description: "Petits fruits poussant en grappes, utilisés pour le vin."},
}

// getFruitByName recherche un fruit par son nom (insensible à la casse)
func getFruitByName(name string) *Fruit {
	for i := range fruits {
		if strings.EqualFold(fruits[i].Name, name) {
			return &fruits[i]
		}
	}
	return nil
}

func main() {
	r := gin.Default()

	// 1. Définir les routes de l'API en premier.
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// GET /api/v1/fruits/:name - Recherche un fruit par nom
			v1.GET("/fruits/:name", func(c *gin.Context) {
				name := c.Param("name")
				fruit := getFruitByName(name)

				if fruit != nil {
					c.JSON(http.StatusOK, fruit)
				} else {
					c.JSON(http.StatusNotFound, gin.H{"error": "Fruit non trouvé"})
				}
			})
		}
	}

	// 2. Ensuite, servir les fichiers statiques du frontend.
	// Le chemin doit être relatif au répertoire depuis lequel le binaire est exécuté.
	// Comme on lance ./backend/app depuis la racine, le chemin correct est ./backend/public.
	r.Use(static.Serve("/", static.LocalFile("./backend/public", true)))

	r.Run(":8080")
}
