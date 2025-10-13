package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gobackend/models"
)


// POST /receipts/upload
// Upload a new receipt
func UploadReceipt(c *gin.Context) {
	// Get the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	// Handle the file upload
	file, err := c.FormFile("receiptImage")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed"})
		return
	}

	// Create a temporary file
	tempDir := "backend/uploads"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.MkdirAll(tempDir, os.ModePerm)
	}
	tempFilePath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temporary file"})
		return
	}
	// Ensure the temporary file is deleted
	defer os.Remove(tempFilePath)

	// Execute the Python OCR script
	cmd := exec.Command("python", "backend/scripts/ocr.py", tempFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("OCR script execution error: %s\n", err)
		log.Printf("OCR script output: %s\n", string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process receipt with OCR"})
		return
	}

	// Parse the JSON output from the script
	var ocrResult struct {
		Montant    float64 `json:"montant"`
		Date       string  `json:"date"`
		Tvac       float64 `json:"tvac"`
		Marchand   string  `json:"marchand"`
		TypeRecu   string  `json:"type_recu"`
		IsAnomaly  bool    `json:"anomalie"`
	}
	if err := json.Unmarshal(output, &ocrResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse OCR results"})
		return
	}

	// Parse date string
	parsedDate, err := time.Parse("2006-01-02", ocrResult.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse date from OCR results"})
		return
	}

	// Create a new receipt model
	newReceipt := models.Receipt{
		UserID:      currentUser.ID,
		Amount:      ocrResult.Montant,
		Date:        parsedDate,
		Vat:         ocrResult.Tvac,
		Merchant:    ocrResult.Marchand,
		ReceiptType: ocrResult.TypeRecu,
		IsAnomaly:   ocrResult.IsAnomaly,
		// In a real app, you would upload to S3/Cloudinary and save the URL here
		FileURL:    tempFilePath,
		RawOcrData: string(output),
	}

	// Save the new receipt to the database
	if result := DB.Create(&newReceipt); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save receipt to database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"receipt": newReceipt})
}

// GET /receipts
// Get all receipts for the current user
func GetReceipts(c *gin.Context) {
	// Get the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	// Find all receipts for the current user
	var receipts []models.Receipt
	if result := DB.Where("user_id = ?", currentUser.ID).Find(&receipts); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve receipts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"receipts": receipts})
}
