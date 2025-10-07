package models

import (
	"time"

	"gorm.io/gorm"
)

// Reservation représente une réservation de place de parking.
type Reservation struct {
	gorm.Model
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	UserID    uint      `json:"user_id"`
	SpotID    uint      `json:"spot_id"`
	TenantID  uint      `json:"tenant_id"`
}
