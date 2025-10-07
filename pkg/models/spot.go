package models

import "gorm.io/gorm"

// Spot représente une place de parking individuelle.
type Spot struct {
	gorm.Model
	SpotNumber string `json:"spot_number"`
	IsOccupied bool   `json:"is_occupied"`
	ParkingID  uint   `json:"parking_id"`
	TenantID   uint   `json:"tenant_id"`
}
