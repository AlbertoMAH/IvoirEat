package handlers

import (
	"errors"
	"net/http"

	"GoBackend/internal/database"
	"GoBackend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is the code for unique_violation
			c.JSON(http.StatusConflict, gin.H{"error": "An account with this email already exists."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}

	// Return the created restaurant
	c.JSON(http.StatusCreated, restaurant)
}

// GetRestaurants handles the API request to list all restaurants
func GetRestaurants(c *gin.Context) {
	var restaurants []models.Restaurant
	// Use Preload("Tables") to also fetch the associated tables for each restaurant
	if err := database.DB.Preload("Tables").Find(&restaurants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurants"})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

// GetRestaurant handles the API request for a single restaurant
func GetRestaurant(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.Restaurant

	if err := database.DB.Preload("Tables").First(&restaurant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}
