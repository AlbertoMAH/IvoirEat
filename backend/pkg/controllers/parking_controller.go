package controllers

import (
	"gobackend/pkg/database"
	"gobackend/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateParking crée un nouveau parking.
func CreateParking(c *gin.Context) {
	var parking models.Parking
	if err := c.ShouldBindJSON(&parking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// The TenantID for a super_admin should be provided in the request body.
	// For tenant_admins, it would be enforced by middleware.
	if err := database.DB.Create(&parking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create parking"})
		return
	}

	c.JSON(http.StatusCreated, parking)
}

// GetParkings récupère tous les parkings.
func GetParkings(c *gin.Context) {
	var parkings []models.Parking
	// The TenantFilterMiddleware is expected to handle filtering for non-super_admin users.
	// A super_admin should get all parkings.
	if err := database.DB.Find(&parkings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve parkings"})
		return
	}

	c.JSON(http.StatusOK, parkings)
}

// GetParking récupère un parking par son ID.
func GetParking(c *gin.Context) {
	var parking models.Parking
	id := c.Param("id")

	if err := database.DB.First(&parking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parking not found"})
		return
	}

	c.JSON(http.StatusOK, parking)
}

// UpdateParking met à jour un parking.
func UpdateParking(c *gin.Context) {
	var parking models.Parking
	id := c.Param("id")

	if err := database.DB.First(&parking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parking not found"})
		return
	}

	// Bind the JSON from the request to a temporary struct to avoid overwriting the ID
	var updatedData models.Parking
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields of the existing parking record
	parking.Name = updatedData.Name
	parking.Location = updatedData.Location
	parking.Capacity = updatedData.Capacity
	parking.TenantID = updatedData.TenantID


	if err := database.DB.Save(&parking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parking"})
		return
	}

	c.JSON(http.StatusOK, parking)
}

// DeleteParking supprime un parking.
func DeleteParking(c *gin.Context) {
	var parking models.Parking
	id := c.Param("id")

	if err := database.DB.First(&parking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parking not found"})
		return
	}

	if err := database.DB.Delete(&parking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete parking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parking deleted successfully"})
}
