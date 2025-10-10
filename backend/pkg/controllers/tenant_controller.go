package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTenant crée un nouveau tenant.
func CreateTenant(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// GetTenants récupère tous les tenants.
func GetTenants(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// GetTenant récupère un tenant par son ID.
func GetTenant(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// UpdateTenant met à jour un tenant.
func UpdateTenant(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// DeleteTenant supprime un tenant.
func DeleteTenant(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}
