package main

import (
	"GoBackend/internal/database"
	"GoBackend/internal/handlers"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	// Set up the router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("internal/templates/*")

	// Redirect root to admin page
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/restaurants")
	})

	// API v1 routes
	api := r.Group("/api/v1")
	{
		api.POST("/restaurants", handlers.CreateRestaurant)
		api.GET("/restaurants", handlers.GetRestaurants)
		api.GET("/restaurants/:id", handlers.GetRestaurant)
		api.GET("/restaurants/:id/availability", handlers.CheckAvailabilityHandler)
	}

	// Admin panel routes
	admin := r.Group("/admin")
	{
		// Restaurant CRUD
		admin.GET("/restaurants", handlers.ListRestaurantsHandler)
		admin.GET("/restaurants/new", handlers.NewRestaurantFormHandler)
		admin.POST("/restaurants", handlers.CreateRestaurantFromFormHandler)
		admin.GET("/restaurants/edit/:id", handlers.EditRestaurantFormHandler)
		admin.POST("/restaurants/edit/:id", handlers.UpdateRestaurantHandler)
		admin.POST("/restaurants/delete/:id", handlers.DeleteRestaurantHandler)

		// Table Management
		admin.GET("/restaurants/:id/tables", handlers.ListTablesHandler)
		admin.POST("/restaurants/:id/tables", handlers.AddTableHandler)
		admin.POST("/tables/status/:id", handlers.UpdateTableStatusHandler)
		admin.POST("/tables/delete/:id", handlers.DeleteTableHandler)

		// Menu of the Day Management
		admin.GET("/restaurants/:id/menu", handlers.MenuOfTheDayHandler)
		admin.POST("/restaurants/:id/menu", handlers.AddDishToMenuHandler)
		admin.GET("/dishes/edit/:id", handlers.EditDishFormHandler)
		admin.POST("/dishes/edit/:id", handlers.UpdateDishHandler)
		admin.POST("/dishes/delete/:id", handlers.DeleteDishFromMenuHandler)
	}

	// Start the server
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
