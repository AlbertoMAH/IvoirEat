package models

import (
	"time"

	"gorm.io/gorm"
)

// DailyDish represents a single dish available on a specific day for a restaurant.
type DailyDish struct {
	gorm.Model
	Name         string    `json:"name"`
	PhotoURL     string    `json:"photo_url"`
	Price        float64   `json:"price"`
	Date         time.Time `json:"date"`
	RestaurantID uint      `json:"restaurant_id" gorm:"index"`
}
