package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// -------------------------
// Models
// -------------------------

// User represents a user of the application
type User struct {
	gorm.Model
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"-"` // Store hashed password, don't expose in JSON
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
	RawOcrData  string    `json:"raw_ocr_data" gorm:"type:jsonb"`
	FileURL     string    `json:"file_url"`
}

// AlertThreshold represents the configurable alert thresholds
type AlertThreshold struct {
	gorm.Model
	HighAmountThreshold float64 `json:"high_amount_threshold"`
	// Add other thresholds here, e.g., duplicate detection time window
}

// -------------------------
// Database
// -------------------------
var DB *gorm.DB

func connectDatabase() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		// Use the one provided by the user if the env var is not set
		dsn = "postgresql://babi_park_user:mORixbz7vPDK9tztKq3jm4ChwsqoPZcp@dpg-d3h57cili9vc73drqmbg-a.frankfurt-postgres.render.com/babi_park"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	err = database.AutoMigrate(&User{}, &Receipt{}, &AlertThreshold{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = database
}

// -------------------------
// API Endpoints
// -------------------------

// ping is a simple health check endpoint
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// -------------------------
// Main
// -------------------------
func main() {
	// Connect to the database
	connectDatabase()

	// Set up the Gin router
	r := gin.Default()

	// A simple ping endpoint to check if the server is running
	r.GET("/ping", ping)


	// TODO: Add auth routes (/auth/register, /auth/login)
	// TODO: Add receipt routes (/receipts)
	// TODO: Add alert routes (/alerts)


	// Start the server
	log.Println("Starting server on port 8080...")
	r.Run(":8080")
}
