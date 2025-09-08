package models

import "gorm.io/gorm"

// ServicePeriod defines a distinct period of service for a restaurant, e.g., Lunch or Dinner.
type ServicePeriod struct {
	gorm.Model

	Name         string `json:"name"`          // e.g., "Lunch", "Dinner"
	OpeningTime  string `json:"opening_time"`  // Format "HH:MM"
	ClosingTime  string `json:"closing_time"`  // Format "HH:MM"
	RestaurantID uint   `json:"restaurant_id"` // Foreign key to Restaurant
}
