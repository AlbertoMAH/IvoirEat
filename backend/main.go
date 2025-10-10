package main

import (
	"gobackend/pkg/database"
	"gobackend/pkg/models"
	"gobackend/pkg/routes"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// StripAppPrefix est un middleware pour supprimer le préfixe /app ajouté par Render.
func StripAppPrefix() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/app") {
			c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/app")
		}
		c.Next()
	}
}

// reverseProxy transfère les requêtes au serveur Next.js.
func reverseProxy(target string) gin.HandlerFunc {
	remoteUrl, err := url.Parse(target)
	if err != nil {
		panic("URL de proxy invalide")
	}
	proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = remoteUrl.Scheme
		req.URL.Host = remoteUrl.Host
		req.Host = remoteUrl.Host
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	database.ConnectDatabase()
	createSuperAdmin()

	router := gin.Default()
	router.Use(StripAppPrefix())

	// Les routes API sont servies par Go.
	routes.SetupRoutes(router)

	// Toutes les autres requêtes sont transférées au serveur Next.js.
	// Le serveur Next.js tournera sur le port 3000 dans le même conteneur.
	nextjsUrl := "http://localhost:3000"
	router.NoRoute(reverseProxy(nextjsUrl))

	router.Run(":8080")
}

func createSuperAdmin() {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		log.Println("ADMIN_EMAIL or ADMIN_PASSWORD not set, skipping super admin creation")
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", adminEmail).First(&user).Error; err == nil {
		log.Println("Super admin user already exists.")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	superAdmin := models.User{
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     "super_admin",
	}

	if err := database.DB.Create(&superAdmin).Error; err != nil {
		log.Fatalf("Failed to create super admin: %v", err)
	}
}
