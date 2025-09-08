package main

import (
	"GoBackend/internal/database"
	"GoBackend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	// Set up the router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("internal/templates/*")

	// API v1 routes
	api := r.Group("/api/v1")
	{
		api.POST("/restaurants", handlers.CreateRestaurant)
	}

	// Admin panel routes
	admin := r.Group("/admin")
	{
		admin.GET("/restaurants", handlers.ListRestaurantsHandler)
		admin.GET("/restaurants/new", handlers.NewRestaurantFormHandler)
		admin.POST("/restaurants", handlers.CreateRestaurantFromFormHandler)
	}

	// Start the server
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
