package handlers

import (
	"net/http"

	"GoBackend/internal/database"
	"GoBackend/internal/models"
	"github.com/gin-gonic/gin"
)

// ListRestaurantsHandler handles the request to display a list of all restaurants.
func ListRestaurantsHandler(c *gin.Context) {
	var restaurants []models.Restaurant
	if err := database.DB.Order("created_at desc").Find(&restaurants).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching restaurants: %v", err)
		return
	}

	// The template will be loaded in main.go
	c.HTML(http.StatusOK, "restaurants.html", gin.H{
		"title":       "Admin - Restaurants",
		"restaurants": restaurants,
	})
}

// NewRestaurantFormHandler displays the form for creating a new restaurant.
func NewRestaurantFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "new_restaurant_form.html", gin.H{
		"title": "Admin - Add New Restaurant",
	})
}

// CreateRestaurantFromFormHandler handles the submission of the new restaurant form.
func CreateRestaurantFromFormHandler(c *gin.Context) {
	// Parse form data and populate the struct
	restaurant := models.Restaurant{
		Name:          c.PostForm("Name"),
		CuisineType:   c.PostForm("CuisineType"),
		Description:   c.PostForm("Description"),
		LogoURL:       c.PostForm("LogoURL"),
		CoverPhotoURL: c.PostForm("CoverPhotoURL"),
		Address:       c.PostForm("Address"),
		GoogleMapsURL: c.PostForm("GoogleMapsURL"),
		Phone:         c.PostForm("Phone"),
		Email:         c.PostForm("Email"),
		WebsiteURL:    c.PostForm("WebsiteURL"),
		HasParking:    c.PostForm("HasParking") == "true",
		HasWifi:       c.PostForm("HasWifi") == "true",
		IsAccessible:  c.PostForm("IsAccessible") == "true",
		IsActive:      c.PostForm("IsActive") == "true",
	}

	// Save to database
	if err := database.DB.Create(&restaurant).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error creating restaurant: %v", err)
		return
	}

	// Redirect to the list page
	c.Redirect(http.StatusFound, "/admin/restaurants")
}
