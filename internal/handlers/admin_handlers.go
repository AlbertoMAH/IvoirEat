package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"GoBackend/internal/database"
	"GoBackend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.String(http.StatusConflict, "Error: The email address '%s' is already in use. Please go back and use a different email.", restaurant.Email)
			return
		}
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
	// Preload the ServicePeriods association
	if err := database.DB.Preload("ServicePeriods").First(&restaurant, id).Error; err != nil {
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

	// Update reservation settings
	durationStr := c.PostForm("AvgReservationDurationInMinutes")
	duration, err := strconv.ParseUint(durationStr, 10, 64)
	if err != nil || duration < 60 {
		// If there is an error or the value is less than 60, set it to 60.
		duration = 60
	}
	restaurant.AvgReservationDurationInMinutes = uint(duration)

	if err := database.DB.Save(&restaurant).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to update restaurant: %v", err)
		return
	}

	// Redirect back to the edit page to show the changes
	c.Redirect(http.StatusFound, "/admin/restaurants/edit/"+id)
}

// --- Service Period Handlers ---

// AddServicePeriodHandler handles adding a new service period to a restaurant.
func AddServicePeriodHandler(c *gin.Context) {
	restaurantID := c.Param("id")
	restID, err := strconv.ParseUint(restaurantID, 10, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid restaurant ID")
		return
	}

	servicePeriod := models.ServicePeriod{
		Name:         c.PostForm("Name"),
		OpeningTime:  c.PostForm("OpeningTime"),
		ClosingTime:  c.PostForm("ClosingTime"),
		RestaurantID: uint(restID),
	}

	if err := database.DB.Create(&servicePeriod).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to create service period: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/restaurants/edit/"+restaurantID)
}

// DeleteServicePeriodHandler handles deleting a service period.
func DeleteServicePeriodHandler(c *gin.Context) {
	periodID := c.Param("id")
	restaurantID := c.Query("restaurant_id") // Needed for the redirect

	if err := database.DB.Delete(&models.ServicePeriod{}, periodID).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete service period: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/restaurants/edit/"+restaurantID)
}

// --- Menu of the Day Handlers ---

// MenuOfTheDayHandler displays the menu for a specific restaurant for the current day.
func MenuOfTheDayHandler(c *gin.Context) {
	restaurantID := c.Param("id")

	var restaurant models.Restaurant
	if err := database.DB.First(&restaurant, restaurantID).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found: %v", err)
		return
	}

	// Get today's date at midnight and tomorrow at midnight to create a date range
	now := time.Now()
	year, month, day := now.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var dishes []models.DailyDish
	if err := database.DB.Where("restaurant_id = ? AND date >= ? AND date < ?", restaurantID, startOfDay, endOfDay).Find(&dishes).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch dishes: %v", err)
		return
	}

	c.HTML(http.StatusOK, "menu_of_the_day.html", gin.H{
		"title":        "Manage Menu of the Day",
		"restaurant":   restaurant,
		"dishes":       dishes,
		"current_date": now.Format("January 02, 2006"),
	})
}

// AddDishToMenuHandler handles adding a new dish to the menu of the day.
func AddDishToMenuHandler(c *gin.Context) {
	restaurantID := c.Param("id")
	name := c.PostForm("Name")
	photoURL := c.PostForm("PhotoURL")
	priceStr := c.PostForm("Price")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid price format: %v", err)
		return
	}

	restID, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid restaurant ID: %v", err)
		return
	}

	dish := models.DailyDish{
		Name:         name,
		PhotoURL:     photoURL,
		Price:        price,
		Date:         time.Now(),
		RestaurantID: uint(restID),
	}

	if err := database.DB.Create(&dish).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to add dish: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/restaurants/"+restaurantID+"/menu")
}

// EditDishFormHandler displays the form for editing a dish.
func EditDishFormHandler(c *gin.Context) {
	dishID := c.Param("id")
	var dish models.DailyDish
	if err := database.DB.First(&dish, dishID).Error; err != nil {
		c.String(http.StatusNotFound, "Dish not found: %v", err)
		return
	}

	c.HTML(http.StatusOK, "edit_dish_form.html", gin.H{
		"title": "Admin - Edit Dish",
		"dish":  dish,
	})
}

// UpdateDishHandler handles the submission of the edit dish form.
func UpdateDishHandler(c *gin.Context) {
	dishID := c.Param("id")
	var dish models.DailyDish
	if err := database.DB.First(&dish, dishID).Error; err != nil {
		c.String(http.StatusNotFound, "Dish not found: %v", err)
		return
	}

	// Update fields from form data
	dish.Name = c.PostForm("Name")
	dish.PhotoURL = c.PostForm("PhotoURL")
	priceStr := c.PostForm("Price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid price format: %v", err)
		return
	}
	dish.Price = price

	// The date is not updated, as it's the "menu of the day"

	if err := database.DB.Save(&dish).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to update dish: %v", err)
		return
	}

	// Convert uint to string for the URL
	restaurantIDStr := strconv.FormatUint(uint64(dish.RestaurantID), 10)
	c.Redirect(http.StatusFound, "/admin/restaurants/"+restaurantIDStr+"/menu")
}

// DeleteDishFromMenuHandler handles deleting a dish from the menu.
func DeleteDishFromMenuHandler(c *gin.Context) {
	dishID := c.Param("id")
	restaurantID := c.Query("restaurant_id") // Get restaurant_id from query param for redirect

	if err := database.DB.Delete(&models.DailyDish{}, dishID).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete dish: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/restaurants/"+restaurantID+"/menu")
}

// --- Table Handlers ---

// ListTablesHandler displays the tables for a specific restaurant.
func ListTablesHandler(c *gin.Context) {
	restaurantID := c.Param("id")

	var restaurant models.Restaurant
	if err := database.DB.First(&restaurant, restaurantID).Error; err != nil {
		c.String(http.StatusNotFound, "Restaurant not found: %v", err)
		return
	}

	var tables []models.Table
	if err := database.DB.Where("restaurant_id = ?", restaurantID).Find(&tables).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch tables: %v", err)
		return
	}

	c.HTML(http.StatusOK, "tables.html", gin.H{
		"title":      "Manage Tables for " + restaurant.Name,
		"restaurant": restaurant,
		"tables":     tables,
	})
}

// AddTableHandler handles the form submission for adding a new table.
func AddTableHandler(c *gin.Context) {
	restaurantID := c.Param("id")
	name := c.PostForm("Name")
	capacityStr := c.PostForm("Capacity")

	capacity, err := strconv.ParseUint(capacityStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid capacity format: %v", err)
		return
	}

	restID, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid restaurant ID: %v", err)
		return
	}

	table := models.Table{
		Name:         name,
		Capacity:     uint(capacity),
		RestaurantID: uint(restID),
	}

	if err := database.DB.Create(&table).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to add table: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/restaurants/"+restaurantID+"/tables")
}

// UpdateTableStatusHandler handles changing the status of a table.
func UpdateTableStatusHandler(c *gin.Context) {
	tableID := c.Param("id")
	newStatus := c.PostForm("Status")
	restaurantID := c.PostForm("restaurant_id")

	var table models.Table
	if err := database.DB.First(&table, tableID).Error; err != nil {
		c.String(http.StatusNotFound, "Table not found: %v", err)
		return
	}

	table.Status = newStatus
	if err := database.DB.Save(&table).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to update table status: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/restaurants/"+restaurantID+"/tables")
}

// DeleteTableHandler handles the deletion of a table.
func DeleteTableHandler(c *gin.Context) {
	tableID := c.Param("id")
	restaurantID := c.Query("restaurant_id") // Get restaurant_id from query param for redirect

	if err := database.DB.Delete(&models.Table{}, tableID).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete table: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/restaurants/"+restaurantID+"/tables")
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
