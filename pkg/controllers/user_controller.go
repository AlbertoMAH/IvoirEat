package controllers

import (
	"net/http"

	"gobackend/pkg/database"
	"gobackend/pkg/models"

	"github.com/gin-gonic/gin"
)

// CreateUser crée un nouvel utilisateur.
func CreateUser(c *gin.Context) {
	// Logique de création d'utilisateur à implémenter.
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// GetUsers récupère tous les utilisateurs.
func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

// GetUser récupère un utilisateur par son ID.
func GetUser(c *gin.Context) {
	// Logique de récupération d'un utilisateur à implémenter.
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// UpdateUser met à jour un utilisateur.
func UpdateUser(c *gin.Context) {
	// Logique de mise à jour d'utilisateur à implémenter.
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// DeleteUser supprime un utilisateur.
func DeleteUser(c *gin.Context) {
	// Logique de suppression d'utilisateur à implémenter.
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}
