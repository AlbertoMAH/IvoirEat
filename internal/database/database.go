package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"GoBackend/internal/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("restaurant.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	log.Println("Database connection successful.")

	err = database.AutoMigrate(&models.Restaurant{}, &models.Table{}, &models.DailyDish{})
	if err != nil {
		log.Fatal("Failed to migrate database!", err)
	}

	log.Println("Database migrated successfully.")

	DB = database
}
