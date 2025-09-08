package main

import (
	"net/http"

	"GoBackend/internal/database"
	"GoBackend/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	// Set up the router
	r := gin.Default()

	// API v1 routes
	api := r.Group("/api/v1")
	{
		api.POST("/restaurants", CreateRestaurant)
		// We can add other restaurant routes here later, e.g., GET, PUT, DELETE
	}

	// Start the server
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}

// CreateRestaurant handles the API request to create a new restaurant
func CreateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant

	// Bind the JSON request body to the restaurant struct
	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the record in the database
	if err := database.DB.Create(&restaurant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}

	// Return the created restaurant
	c.JSON(http.StatusCreated, restaurant)
}
