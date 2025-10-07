package models

import "gorm.io/gorm"

// Tenant représente une organisation ou un client dans le système.
type Tenant struct {
	gorm.Model
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}
