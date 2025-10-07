package models

import (
	"gorm.io/gorm"
)

// Log représente un enregistrement d'activité dans le système.
type Log struct {
	gorm.Model
	Action   string `json:"action"` // Ex: "create_user", "delete_parking", etc.
	UserID   uint   `json:"user_id"`
	TenantID uint   `json:"tenant_id"`
}
