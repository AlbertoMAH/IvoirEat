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

// EditRestaurantFormHandler displays the form for editing an existing restaurant.
func EditRestaurantFormHandler(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.Restaurant
	if err := database.DB.First(&restaurant, id).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found: %v", err)
		return
	}

	c.HTML(http.StatusOK, "edit_restaurant_form.html", gin.H{
		"title":      "Admin - Edit Restaurant",
		"restaurant": restaurant,
	})
}

// UpdateRestaurantHandler handles the submission of the edit restaurant form.
func UpdateRestaurantHandler(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.Restaurant
	if err := database.DB.First(&restaurant, id).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found: %v", err)
		return
	}

	// Update fields from form data
	restaurant.Name = c.PostForm("Name")
	restaurant.CuisineType = c.PostForm("CuisineType")
	restaurant.Description = c.PostForm("Description")
	restaurant.LogoURL = c.PostForm("LogoURL")
	restaurant.CoverPhotoURL = c.PostForm("CoverPhotoURL")
	restaurant.Address = c.PostForm("Address")
	restaurant.GoogleMapsURL = c.PostForm("GoogleMapsURL")
	restaurant.Phone = c.PostForm("Phone")
	restaurant.Email = c.PostForm("Email")
	restaurant.WebsiteURL = c.PostForm("WebsiteURL")
	restaurant.HasParking = c.PostForm("HasParking") == "true"
	restaurant.HasWifi = c.PostForm("HasWifi") == "true"
	restaurant.IsAccessible = c.PostForm("IsAccessible") == "true"
	restaurant.IsActive = c.PostForm("IsActive") == "true"

	if err := database.DB.Save(&restaurant).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to update restaurant: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/restaurants")
}

// DeleteRestaurantHandler handles the deletion of a restaurant.
func DeleteRestaurantHandler(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Restaurant{}, id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete restaurant: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/restaurants")
}
