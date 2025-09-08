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

	// Get today's date range
	now := time.Now()
	year, month, day := now.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Preload associations
	if err := database.DB.
		Preload("Tables").
		Preload("ServicePeriods").
		Preload("DailyDishes", "date >= ? AND date < ?", startOfDay, endOfDay).
		Find(&restaurants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurants"})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

// GetRestaurant handles the API request for a single restaurant
func GetRestaurant(c *gin.Context) {
	id := c.Param("id")
	var restaurant models.Restaurant

	// Get today's date range
	now := time.Now()
	year, month, day := now.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := database.DB.
		Preload("Tables").
		Preload("ServicePeriods").
		Preload("DailyDishes", "date >= ? AND date < ?", startOfDay, endOfDay).
		First(&restaurant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// CheckAvailabilityHandler handles the API request to check for available reservation slots based on total capacity.
func CheckAvailabilityHandler(c *gin.Context) {
	restaurantID := c.Param("id")
	dateStr := c.Query("date")       // e.g., "2024-12-25"
	partySizeStr := c.Query("party_size") // e.g., "4"

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

	// --- 2. Fetch Restaurant, Tables, and Service Periods ---
	var restaurant models.Restaurant
	if err := database.DB.Preload("Tables").Preload("ServicePeriods").First(&restaurant, restaurantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}
	if len(restaurant.ServicePeriods) == 0 {
		c.JSON(http.StatusOK, gin.H{"available_slots": []string{}, "message": "This restaurant is not configured for reservations."})
		return
	}

	// --- 3. Calculate Total Capacity ---
	var totalCapacity int
	for _, table := range restaurant.Tables {
		if table.Status == "available" {
			totalCapacity += int(table.Capacity)
		}
	}
	if totalCapacity < partySize {
		c.JSON(http.StatusOK, gin.H{"available_slots": []string{}, "message": "The restaurant does not have enough total capacity for a party of this size."})
		return
	}

	// --- 4. Generate Potential Time Slots & Check Availability ---
	reservationDuration := time.Duration(restaurant.AvgReservationDurationInMinutes) * time.Minute
	timeSlotIncrement := 15 * time.Minute
	availableSlots := make(map[string]bool)

	// Iterate over each service period for the restaurant
	for _, period := range restaurant.ServicePeriods {
		openingTime, _ := time.Parse("15:04", period.OpeningTime)
		closingTime, _ := time.Parse("15:04", period.ClosingTime)

		// Start from the service's opening time on the requested date
		currentTimeSlot := time.Date(date.Year(), date.Month(), date.Day(), openingTime.Hour(), openingTime.Minute(), 0, 0, time.UTC)
		// The last possible time to start a reservation
		lastBookingTime := time.Date(date.Year(), date.Month(), date.Day(), closingTime.Hour(), closingTime.Minute(), 0, 0, time.UTC).Add(-reservationDuration)

		for currentTimeSlot.Before(lastBookingTime) || currentTimeSlot.Equal(lastBookingTime) {
			slotStartTime := currentTimeSlot
			slotEndTime := slotStartTime.Add(reservationDuration)

			// --- Calculate booked capacity for this specific time slot ---
			var bookedCapacity int
			var overlappingReservations []models.Reservation

			// Find all table IDs for the current restaurant
			var tableIDs []uint
			for _, table := range restaurant.Tables {
				tableIDs = append(tableIDs, table.ID)
			}

			// Find reservations that overlap with the current time slot across all tables of the restaurant
			database.DB.Model(&models.Reservation{}).
				Where("table_id IN (?)", tableIDs).
				Where("reservation_start_time < ? AND reservation_end_time > ?", slotEndTime, slotStartTime).
				Find(&overlappingReservations)

			for _, r := range overlappingReservations {
				bookedCapacity += int(r.NumberOfGuests)
			}

			// Check if there is enough capacity
			if (totalCapacity - bookedCapacity) >= partySize {
				slotStr := slotStartTime.Format("15:04")
				availableSlots[slotStr] = true
			}

			currentTimeSlot = currentTimeSlot.Add(timeSlotIncrement)
		}
	}

	// --- 5. Format and Return Response ---
	finalSlots := []string{}
	for slot := range availableSlots {
		finalSlots = append(finalSlots, slot)
	}

	c.JSON(http.StatusOK, gin.H{"available_slots": finalSlots})
}
