package controllers

import (
	"gobackend/pkg/database"
	"gobackend/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats récupère les statistiques pour le tableau de bord du super admin.
func GetDashboardStats(c *gin.Context) {
	var totalParkings int64
	if err := database.DB.Model(&models.Parking{}).Count(&totalParkings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count parkings"})
		return
	}

	var totalSpots int64
	if err := database.DB.Model(&models.Spot{}).Count(&totalSpots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count spots"})
		return
	}

	var availableSpots int64
	if err := database.DB.Model(&models.Spot{}).Where("is_occupied = ?", false).Count(&availableSpots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count available spots"})
		return
	}

	var totalUsers int64
	if err := database.DB.Model(&models.User{}).Where("role != ?", "super_admin").Count(&totalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	var todayReservations int64
	if err := database.DB.Model(&models.Reservation{}).Where("start_time >= ?", today).Count(&todayReservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count today's reservations"})
		return
	}

	// Récupérer les données pour le graphique des 7 derniers jours
	var weeklyReservations []struct {
		Date  time.Time `json:"date"`
		Count int       `json:"count"`
	}
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	if err := database.DB.Model(&models.Reservation{}).
		Select("DATE(start_time) as date, COUNT(*) as count").
		Where("start_time >= ?", sevenDaysAgo).
		Group("DATE(start_time)").
		Order("date ASC").
		Scan(&weeklyReservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weekly reservations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_parkings":      totalParkings,
		"total_spots":         totalSpots,
		"available_spots":     availableSpots,
		"total_users":         totalUsers,
		"today_reservations":  todayReservations,
		"weekly_reservations": weeklyReservations,
	})
}
