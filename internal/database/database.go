package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"GoBackend/internal/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	// Check if we should use SQLite
	if os.Getenv("USE_SQLITE") == "true" {
		DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to SQLite database! ", err)
		}
		log.Println("SQLite database connection successful (in-memory).")
	} else {
		// Default to PostgreSQL
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("DATABASE_URL environment variable is not set")
		}
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL database! ", err)
		}
		log.Println("PostgreSQL database connection successful.")
	}

	// Migrate the schema
	err = DB.AutoMigrate(&models.Restaurant{}, &models.Table{}, &models.DailyDish{}, &models.Reservation{}, &models.ServicePeriod{})
	if err != nil {
		log.Fatal("Failed to migrate database!", err)
	}

	log.Println("Database migrated successfully.")
}
