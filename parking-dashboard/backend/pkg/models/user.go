package models

import (
	"gorm.io/gorm"
)

// User représente un utilisateur dans le système.
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role"` // Ex: "super_admin", "tenant_admin", "parking_agent", "user"
	TenantID uint   `json:"tenant_id"`
}
