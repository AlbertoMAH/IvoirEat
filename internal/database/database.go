package database

import (
	"log"
	// "os" // No longer needed

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"GoBackend/internal/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Hardcode the DSN to ensure the correct one is used, including SSL mode.
	dsn := "postgresql://ivoireat_db_user:KpkBYEwrBPpi3EHfgF7Hp6cTswEO4Jhi@dpg-d2vcuov5r7bs73co8bmg-a.frankfurt-postgres.render.com/ivoireat_db?sslmode=require"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database! ", err)
	}

	log.Println("Database connection successful.")

	err = database.AutoMigrate(&models.Restaurant{}, &models.Table{}, &models.DailyDish{}, &models.Reservation{}, &models.ServicePeriod{})
	if err != nil {
		log.Fatal("Failed to migrate database!", err)
	}

	log.Println("Database migrated successfully.")

	DB = database
}
