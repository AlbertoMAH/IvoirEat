package models

import (
	"time"

	"gorm.io/gorm"
)

// Reservation holds the details of a booking.
type Reservation struct {
	gorm.Model
	CustomerName         string    `json:"customer_name"`
	CustomerPhone        string    `json:"customer_phone"`
	CustomerEmail        string    `json:"customer_email"`
	NumberOfGuests       uint      `json:"number_of_guests"`
	ReservationStartTime time.Time `json:"reservation_start_time"`
	ReservationEndTime   time.Time `json:"reservation_end_time"`
	Status               string    `json:"status" gorm:"default:'confirmed'"` // e.g., "confirmed", "cancelled", "completed"
	RestaurantID         uint      `json:"restaurant_id" gorm:"index"`
	TableID              uint      `json:"table_id" gorm:"index"`
}
