package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gobackend/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "postgresql://babi_park_user:mORixbz7vPDK9tztKq3jm4ChwsqoPZcp@dpg-d3h57cili9vc73drqmbg-a.frankfurt-postgres.render.com/babi_park"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	err = database.AutoMigrate(&models.User{}, &models.Receipt{}, &models.AlertThreshold{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = database
}
