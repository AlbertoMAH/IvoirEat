package models

import "gorm.io/gorm"

// Table represents a single table within a restaurant.
type Table struct {
	gorm.Model
	Name         string `json:"name"`
	Capacity     uint   `json:"capacity"`
	Status       string `json:"status" gorm:"default:'available'"` // e.g., "available", "occupied", "out_of_service"
	RestaurantID uint   `json:"restaurant_id" gorm:"index"`      // Foreign key for Restaurant
}
