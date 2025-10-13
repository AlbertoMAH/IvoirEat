package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gobackend/controllers"
	"gobackend/database"
	"gobackend/middleware"
)

// ping is a simple health check endpoint
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func main() {
	// Connect to the database
	database.ConnectDatabase()

	// Pass the database connection to the controllers
	controllers.DB = database.DB

	// Set up the Gin router
	r := gin.Default()

	// Group public routes
	public := r.Group("/api")
	{
		public.GET("/ping", ping)
	}

	// Group auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Group protected routes
	protected := r.Group("/api")
	protected.Use(middleware.RequireAuth)
	{
		// This is a test route to check if the auth middleware works
		protected.GET("/validate", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(http.StatusOK, gin.H{"message": "I'm logged in", "user": user})
		})
		// TODO: Add other protected routes for receipts, alerts, etc.
	}

	// Start the server
	log.Println("Starting server on port 8080...")
	r.Run(":8080")
}
