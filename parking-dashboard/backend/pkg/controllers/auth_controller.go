package controllers

import (
	"gobackend/pkg/database"
	"gobackend/pkg/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginInput définit la structure pour le corps de la requête de connexion.
type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login gère l'authentification des utilisateurs.
func Login(c *gin.Context) {
	var input LoginInput

	// Valider le corps de la requête
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trouver l'utilisateur par email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe invalide"})
		return
	}

	// Comparer le mot de passe fourni avec le hash stocké
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe invalide"})
		return
	}

	// Générer le token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(), // Le token expire dans 30 jours
	})

	// Signer le token avec la clé secrète
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "La clé secrète JWT n'est pas configurée"})
		return
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de générer le token"})
		return
	}

	// Renvoyer le token
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
