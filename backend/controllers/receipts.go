package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/gin-gonic/gin"
	"gobackend/models"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	G  *genkit.Genkit
)

// ReceiptData defines the structure for the data we want the AI to extract.
type ReceiptData struct {
	Merchant    string  `json:"merchant"`
	Total       float64 `json:"total"`
	Date        string  `json:"date"` // Expecting "YYYY-MM-DD"
	Category    string  `json:"category"`
	Description string  `json:"description"`
}

func UploadReceipt(c *gin.Context) {
	// 1. Get user and file from the request
	user, _ := c.Get("user")
	currentUser := user.(models.User)
	fileHeader, err := c.FormFile("receiptImage")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image"})
		return
	}
	defer file.Close()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image content"})
		return
	}
	mimeType := fileHeader.Header.Get("Content-Type")

	// 2. Define the multimodal prompt for Genkit
	prompt := []*ai.Part{
		ai.NewTextPart("Extract the following information from the receipt image and return it as a valid JSON object. " +
			"The fields are: merchant (string), total (float64), date (string, as YYYY-MM-DD), " +
			"category (string, one of: Transport, Restaurant, Hébergement, Achats, Services, Santé, Divers), " +
			"and description (string, a brief summary of the items)."),
		ai.NewDataPart(fileContents, mimeType),
	}

	// 3. Run the generation with the prompt
	resp, err := genkit.Generate(c.Request.Context(), G,
		ai.WithModelName("googleai/gemini-2.5-flash"),
		ai.WithPrompt(prompt),
		ai.WithOutputFormat(ai.OutputFormatJSON),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to extract data from receipt using AI",
			"details": err.Error(),
		})
		return
	}

	// 4. Extract and parse the JSON response from the model
	var extractedData ReceiptData
	jsonText := resp.Text()
	if err := json.Unmarshal([]byte(jsonText), &extractedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse AI response",
			"details": err.Error(),
		})
		return
	}

	// 5. Parse the date and create the database model
	parsedDate, err := time.Parse("2006-01-02", extractedData.Date)
	if err != nil {
		parsedDate = time.Now()
	}

	// 6. Save the new receipt to the database
	newReceipt := models.Receipt{
		UserID:      currentUser.ID,
		Amount:      extractedData.Total,
		Date:        parsedDate,
		Vat:         0.0,
		Merchant:    extractedData.Merchant,
		ReceiptType: extractedData.Category,
		Description: extractedData.Description,
		IsAnomaly:   false,
		FileURL:     "",
		RawOcrData:  jsonText, // Store the raw JSON from the AI
	}

	if result := DB.Create(&newReceipt); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save receipt to database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"receipt": newReceipt})
}

// GET /receipts
func GetReceipts(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(models.User)
	var receipts []models.Receipt
	DB.Where("user_id = ?", currentUser.ID).Order("date desc").Find(&receipts)
	c.JSON(http.StatusOK, gin.H{"receipts": receipts})
}

// GET /receipts/:id
func GetReceiptByID(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(models.User)
	var receipt models.Receipt
	receiptID := c.Param("id")

	if err := DB.Where("id = ? AND user_id = ?", receiptID, currentUser.ID).First(&receipt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"receipt": receipt})
}

// DELETE /receipts/:id
func DeleteReceipt(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(models.User)
	receiptID := c.Param("id")

	var receipt models.Receipt
	if err := DB.Where("id = ? AND user_id = ?", receiptID, currentUser.ID).First(&receipt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	if err := DB.Delete(&receipt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete receipt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Receipt deleted successfully"})
}

// PUT /receipts/:id
func UpdateReceipt(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(models.User)
	receiptID := c.Param("id")

	var receipt models.Receipt
	if err := DB.Where("id = ? AND user_id = ?", receiptID, currentUser.ID).First(&receipt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	var input models.Receipt
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Model(&receipt).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receipt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"receipt": receipt})
}
