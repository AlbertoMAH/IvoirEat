package database

import (
	"log"
	"os"

	"gobackend/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL != "" {
		// Production: PostgreSQL
		log.Println("Connecting to PostgreSQL database...")
		DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	} else {
		// Development: SQLite
		log.Println("Connecting to SQLite database...")
		DB, err = gorm.Open(sqlite.Open("restaurant.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Migration des modèles
	err = DB.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.Parking{},
		&models.Spot{},
		&models.Reservation{},
		&models.Log{},
	)
	if err != nil {
		log.Fatal("Erreur de migration :", err)
	}
}
