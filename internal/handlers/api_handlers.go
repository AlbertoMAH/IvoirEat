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

// CheckAvailabilityHandler handles the API request to check for available reservation slots.
func CheckAvailabilityHandler(c *gin.Context) {
	restaurantID := c.Param("id")
	dateStr := c.Query("date") // e.g., "2024-12-25"
	partySizeStr := c.Query("party_size")

	// --- 1. Parse and Validate Inputs ---
	if dateStr == "" || partySizeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date and party_size query parameters are required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Please use YYYY-MM-DD."})
		return
	}

	partySize, err := strconv.Atoi(partySizeStr)
	if err != nil || partySize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid party_size. Must be a positive integer."})
		return
	}

	// --- 2. Fetch Restaurant and its Rules ---
	var restaurant models.Restaurant
	if err := database.DB.First(&restaurant, restaurantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	// --- 3. Fetch Suitable Tables ---
	var suitableTables []models.Table
	if err := database.DB.Where("restaurant_id = ? AND capacity >= ?", restaurantID, partySize).Find(&suitableTables).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tables"})
		return
	}
	if len(suitableTables) == 0 {
		c.JSON(http.StatusOK, gin.H{"available_slots": []string{}, "message": "No tables large enough for the requested party size."})
		return
	}

	// --- 4. Generate Potential Time Slots & Check Availability ---
	openingTime, _ := time.Parse("15:04", restaurant.OpeningTime)
	closingTime, _ := time.Parse("15:04", restaurant.ClosingTime)
	reservationDuration := time.Duration(restaurant.AvgReservationDurationInMinutes) * time.Minute
	timeSlotIncrement := 15 * time.Minute // Check every 15 minutes

	availableSlots := make(map[string]bool)

	// Start from the restaurant's opening time on the requested date
	currentTimeSlot := time.Date(date.Year(), date.Month(), date.Day(), openingTime.Hour(), openingTime.Minute(), 0, 0, time.UTC)

	// The last possible time to start a reservation
	lastBookingTime := time.Date(date.Year(), date.Month(), date.Day(), closingTime.Hour(), closingTime.Minute(), 0, 0, time.UTC).Add(-reservationDuration)

	for currentTimeSlot.Before(lastBookingTime) || currentTimeSlot.Equal(lastBookingTime) {
		slotStartTime := currentTimeSlot
		slotEndTime := slotStartTime.Add(reservationDuration)

		// Check if this slot is available on ANY of the suitable tables
		for _, table := range suitableTables {
			var overlappingReservations int64
			database.DB.Model(&models.Reservation{}).
				Where("table_id = ?", table.ID).
				Where("reservation_start_time < ? AND reservation_end_time > ?", slotEndTime, slotStartTime).
				Count(&overlappingReservations)

			if overlappingReservations == 0 {
				// This slot is available for this table
				slotStr := slotStartTime.Format("15:04")
				availableSlots[slotStr] = true
				break // No need to check other tables for this same time slot
			}
		}
		currentTimeSlot = currentTimeSlot.Add(timeSlotIncrement)
	}

	// --- 5. Format and Return Response ---
	finalSlots := []string{}
	for slot := range availableSlots {
		finalSlots = append(finalSlots, slot)
	}

	c.JSON(http.StatusOK, gin.H{"available_slots": finalSlots})
}
