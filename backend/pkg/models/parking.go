package models

import "gorm.io/gorm"

// Parking représente un parking géré par un tenant.
type Parking struct {
	gorm.Model
	Name     string `json:"name"`
	Location string `json:"location"`
	Capacity int    `json:"capacity"`
	TenantID uint   `json:"tenant_id"`
}
