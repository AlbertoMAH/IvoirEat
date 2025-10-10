package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:public
var staticFiles embed.FS

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

	// 2. Servir les fichiers statiques du frontend directement depuis le binaire.
	// Crée un sous-système de fichiers qui pointe vers le contenu de notre dossier 'public' embarqué.
	subFS, err := fs.Sub(staticFiles, "public")
	if err != nil {
		panic("Erreur: impossible de créer le sous-système de fichiers pour les fichiers statiques: " + err.Error())
	}

	// Utilise NoRoute pour intercepter toutes les requêtes qui ne correspondent pas à une route API.
	// C'est la méthode idéale pour servir une Single Page Application (SPA).
	r.NoRoute(func(c *gin.Context) {
		// Récupère le chemin du fichier demandé depuis l'URL.
		filePath := strings.TrimPrefix(c.Request.URL.Path, "/")

		// Tente d'ouvrir le fichier pour vérifier s'il existe dans le FS embarqué.
		file, err := subFS.Open(filePath)
		if err != nil {
			// Si le fichier n'existe pas (c'est une route gérée par le frontend),
			// on sert 'index.html', qui est le point d'entrée de la SPA.
			c.FileFromFS("index.html", http.FS(subFS))
			return
		}
		file.Close() // Il est important de fermer le fichier après vérification.

		// Si le fichier existe, on le sert.
		c.FileFromFS(filePath, http.FS(subFS))
	})


	r.Run(":8080")
}
