package database

import (
	"log"

	"gobackend/pkg/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("restaurant.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Erreur connexion base :", err)
	}

	// Migration des modèles
	err = database.AutoMigrate(
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

	DB = database
}
