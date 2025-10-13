package controllers

import (
	"net/http"
	"gobackend/models"
	"github.com/gin-gonic/gin"
)

// GET /alerts/thresholds
// Get the current alert thresholds
func GetThresholds(c *gin.Context) {
	var threshold models.AlertThreshold
	// We'll just get the first one for now, assuming there's only one row for thresholds
	if result := DB.First(&threshold); result.Error != nil {
		// If no threshold is set, create a default one
		if result.Error.Error() == "record not found" {
			defaultThreshold := models.AlertThreshold{HighAmountThreshold: 100.0} // Default value
			DB.Create(&defaultThreshold)
			c.JSON(http.StatusOK, gin.H{"thresholds": defaultThreshold})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve thresholds"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"thresholds": threshold})
}

// POST /alerts/thresholds
// Update the alert thresholds
func UpdateThresholds(c *gin.Context) {
	var threshold models.AlertThreshold
	// Get the first threshold record
	if result := DB.First(&threshold); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Thresholds not found"})
		return
	}

	// Bind the request body to update the threshold
	if err := c.ShouldBindJSON(&threshold); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated threshold
	if result := DB.Save(&threshold); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update thresholds"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"thresholds": threshold})
}
