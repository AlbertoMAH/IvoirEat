package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateSpot crée une nouvelle place de parking.
func CreateSpot(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// GetSpots récupère toutes les places de parking.
func GetSpots(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// GetSpot récupère une place de parking par son ID.
func GetSpot(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// UpdateSpot met à jour une place de parking.
func UpdateSpot(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// DeleteSpot supprime une place de parking.
func DeleteSpot(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}
