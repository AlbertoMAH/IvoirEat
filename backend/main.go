package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
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
	// Initialize Genkit, expecting GEMINI_API_KEY to be in the environment
	ctx := context.Background()
	if os.Getenv("GEMINI_API_KEY") == "" {
		// In a real app, you might not want to log.Fatal here,
		// but for this project, it's a clear indicator of a missing config.
		log.Fatal("GEMINI_API_KEY environment variable not set")
	}

	// Initialize Genkit with the Google AI plugin
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
	)
	controllers.G = g // Make genkit instance available to controllers

	// Connect to the database
	database.ConnectDatabase()

	// Pass the database connection to the controllers
	controllers.DB = database.DB

	// Set up the Gin router
	r := gin.Default()

	// All API routes will be under the /api group
	api := r.Group("/api")
	{
		// Public routes
		api.GET("/ping", ping)

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// Protected routes - everything in here requires authentication
		protected := api.Group("/")
		protected.Use(middleware.RequireAuth)
		{
			// This is a test route to check if the auth middleware works
			protected.GET("/validate", func(c *gin.Context) {
				user, _ := c.Get("user")
				c.JSON(http.StatusOK, gin.H{"message": "I'm logged in", "user": user})
			})

			// Receipts routes
			receipts := protected.Group("/receipts")
			{
				receipts.POST("/upload", controllers.UploadReceipt)
				receipts.GET("", controllers.GetReceipts) // Corrected route path
			}

			// Alerts routes
			alerts := protected.Group("/alerts")
			{
				alerts.GET("/thresholds", controllers.GetThresholds)
				alerts.POST("/thresholds", controllers.UpdateThresholds)
			}
		}
	}

	// Start the server
	log.Println("Starting server on port 8080...")
	r.Run(":8080")
}
