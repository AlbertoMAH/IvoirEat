package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"GoBackend/internal/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database! ", err)
	}

	log.Println("Database connection successful.")

	err = database.AutoMigrate(&models.Restaurant{}, &models.Table{}, &models.DailyDish{})
	if err != nil {
		log.Fatal("Failed to migrate database!", err)
	}

	log.Println("Database migrated successfully.")

	DB = database
}
