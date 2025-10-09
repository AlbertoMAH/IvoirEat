package routes

import (
	"gobackend/pkg/controllers"
	"gobackend/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configure toutes les routes de l'application.
func SetupRoutes(router *gin.Engine) {
	// Route de health check à la racine
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Welcome to the Parking Reservation API",
		})
	})

	api := router.Group("/api")
	{
		// Routes publiques (ex: login, register)
		api.POST("/login", controllers.Login)

		// Routes protégées
		v1 := api.Group("/v1")
		v1.Use(middlewares.JWTMiddleware())
		v1.Use(middlewares.TenantFilterMiddleware())
		{
			// Routes pour les utilisateurs
			users := v1.Group("/users")
			{
				users.POST("/", controllers.CreateUser)
				users.GET("/", controllers.GetUsers)
				users.GET("/:id", controllers.GetUser)
				users.PUT("/:id", controllers.UpdateUser)
				users.DELETE("/:id", controllers.DeleteUser)
			}

			// Routes pour les tenants
			tenants := v1.Group("/tenants")
			{
				tenants.POST("/", controllers.CreateTenant)
				tenants.GET("/", controllers.GetTenants)
				tenants.GET("/:id", controllers.GetTenant)
				tenants.PUT("/:id", controllers.UpdateTenant)
				tenants.DELETE("/:id", controllers.DeleteTenant)
			}

			// Routes pour les parkings
			parkings := v1.Group("/parkings")
			{
				parkings.POST("/", controllers.CreateParking)
				parkings.GET("/", controllers.GetParkings)
				parkings.GET("/:id", controllers.GetParking)
				parkings.PUT("/:id", controllers.UpdateParking)
				parkings.DELETE("/:id", controllers.DeleteParking)
			}

			// Routes pour les spots
			spots := v1.Group("/spots")
			{
				spots.POST("/", controllers.CreateSpot)
				spots.GET("/", controllers.GetSpots)
				spots.GET("/:id", controllers.GetSpot)
				spots.PUT("/:id", controllers.UpdateSpot)
				spots.DELETE("/:id", controllers.DeleteSpot)
			}

			// Routes pour les réservations
			reservations := v1.Group("/reservations")
			{
				reservations.POST("/", controllers.CreateReservation)
				reservations.GET("/", controllers.GetReservations)
				reservations.GET("/:id", controllers.GetReservation)
				reservations.PUT("/:id", controllers.UpdateReservation)
				reservations.DELETE("/:id", controllers.DeleteReservation)
			}

			// Route pour les statistiques du dashboard
			dashboard := v1.Group("/dashboard")
			{
				dashboard.GET("/stats", controllers.GetDashboardStats)
			}
		}
	}
}
