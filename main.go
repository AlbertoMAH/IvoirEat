package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// -------------------------
// Modèles
// -------------------------
type Admin struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type Restaurant struct {
	gorm.Model
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Tables      int    `json:"tables"`
	Capacity    int    `json:"capacity"`
	AdminID     uint   `json:"admin_id"`
}

// -------------------------
// Base de données
// -------------------------
var DB *gorm.DB

func connectDatabase() {
	database, err := gorm.Open(sqlite.Open("restaurant.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Erreur connexion base :", err)
	}
	database.AutoMigrate(&Admin{}, &Restaurant{})
	DB = database
}

// -------------------------
// Input inscription
// -------------------------
type RestaurantRegisterInput struct {
	Admin struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
	} `json:"admin"`
	Restaurant struct {
		Name        string `json:"name"`
		LogoURL     string `json:"logo_url"`
		Description string `json:"description"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Tables      int    `json:"tables"`
		Capacity    int    `json:"capacity"`
	} `json:"restaurant"`
}

// -------------------------
// Routes
// -------------------------
func registerRestaurant(c *gin.Context) {
	var input RestaurantRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin := Admin{
		Name:     input.Admin.Name,
		Email:    input.Admin.Email,
		Password: input.Admin.Password, // ⚠️ hash à ajouter plus tard
		Phone:    input.Admin.Phone,
	}
	DB.Create(&admin)

	restaurant := Restaurant{
		Name:        input.Restaurant.Name,
		LogoURL:     input.Restaurant.LogoURL,
		Description: input.Restaurant.Description,
		Address:     input.Restaurant.Address,
		Phone:       input.Restaurant.Phone,
		Email:       input.Restaurant.Email,
		Tables:      input.Restaurant.Tables,
		Capacity:    input.Restaurant.Capacity,
		AdminID:     admin.ID,
	}
	DB.Create(&restaurant)

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Restaurant inscrit avec succès",
		"restaurant": restaurant,
	})
}

// Page HTML Dashboard Admin
func adminDashboard(c *gin.Context) {
	var restaurants []Restaurant
	DB.Find(&restaurants)

	tmpl := template.Must(template.New("dashboard").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Dashboard Admin</title>
	<style>
		body { font-family: Arial; margin: 20px; }
		table { border-collapse: collapse; width: 100%; }
		th, td { border: 1px solid #ccc; padding: 8px; text-align: left; }
		th { background: #eee; }
	</style>
</head>
<body>
	<h1>Dashboard Admin</h1>
	<h2>Restaurants</h2>
	<table>
		<tr>
			<th>ID</th><th>Nom</th><th>Adresse</th><th>Tables</th><th>Capacité</th>
		</tr>
		{{range .}}
		<tr>
			<td>{{.ID}}</td>
			<td>{{.Name}}</td>
			<td>{{.Address}}</td>
			<td>{{.Tables}}</td>
			<td>{{.Capacity}}</td>
		</tr>
		{{end}}
	</table>
	<p>Pour ajouter un restaurant, utilisez l’API POST /register</p>
</body>
</html>
	`))

	tmpl.Execute(c.Writer, restaurants)
}

// Ping API endpoint for health check
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// -------------------------
// Main
// -------------------------
func main() {
	connectDatabase()

	r := gin.Default()
	r.POST("/register", registerRestaurant)
	r.GET("/admin", adminDashboard)
	r.GET("/ping", ping)

	r.Run(":8080")
}

