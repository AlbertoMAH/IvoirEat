package models

import "gorm.io/gorm"

type Restaurant struct {
    gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt

    // Informations de base
    Name          string `json:"name"`
    CuisineType   string `json:"cuisine_type"`
    Description   string `json:"description"`
    LogoURL       string `json:"logo_url"`
    CoverPhotoURL string `json:"cover_photo_url"`

    // Localisation
    Address       string `json:"address"`
    GoogleMapsURL string `json:"google_maps_url"`

    // Coordonnées de contact
    Phone        string `json:"phone"`
    Email        string `json:"email" gorm:"unique"`
    WebsiteURL   string `json:"website_url"`

    // Services & Options
    HasParking   bool `json:"has_parking"`
    HasWifi      bool `json:"has_wifi"`
    IsAccessible bool `json:"is_accessible"`

    // Statut
    IsActive bool `json:"is_active" gorm:"default:true"`

    // Reservation Settings
    AvgReservationDurationInMinutes uint `json:"avg_reservation_duration_in_minutes" gorm:"default:60"`
    SlotIntervalInMinutes           uint `json:"slot_interval_in_minutes" gorm:"default:60"`

    // Associations
    Tables         []Table         `json:"tables"`          // A restaurant has many tables
    DailyDishes    []DailyDish     `json:"daily_dishes"`    // A restaurant has many daily dishes
    ServicePeriods []ServicePeriod `json:"service_periods"` // A restaurant has many service periods
}
