package models

import (
	"time"
	"gorm.io/gorm"
)

// User represents a user of the application
type User struct {
	gorm.Model
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"-"`
	Receipts []Receipt `json:"receipts"`
}

// Receipt represents a single expense receipt
type Receipt struct {
	gorm.Model
	UserID      uint      `json:"user_id"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Vat         float64   `json:"vat"`
	Merchant    string    `json:"merchant"`
	ReceiptType string    `json:"receipt_type"`
	IsAnomaly   bool      `json:"is_anomaly"`
	RawOcrData  string    `json:"raw_ocr_data" gorm:"type:text"`
	FileURL     string    `json:"file_url"`
}

// AlertThreshold represents the configurable alert thresholds
type AlertThreshold struct {
	gorm.Model
	HighAmountThreshold float64 `json:"high_amount_threshold"`
}
