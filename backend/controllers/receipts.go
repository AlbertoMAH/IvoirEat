package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gobackend/models"
)


// POST /receipts/upload
// Upload a new receipt by calling the OCR microservice
func UploadReceipt(c *gin.Context) {
	// 1. Get the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	// 2. Handle the file upload from the request
	fileHeader, err := c.FormFile("receiptImage")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed"})
		return
	}

	// 3. Prepare the file for the HTTP request to the OCR service
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer file.Close()

	// Create a buffer to store our request body
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form field for the file
	part, err := writer.CreateFormFile("receiptImage", fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create form file"})
		return
	}

	// Copy the file content to the form field
	_, err = io.Copy(part, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file content"})
		return
	}
	writer.Close()

	// 4. Make the HTTP POST request to the OCR service
	ocrServiceURL := os.Getenv("OCR_SERVICE_URL")
	if ocrServiceURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OCR service URL not configured"})
		return
	}

	req, err := http.NewRequest("POST", ocrServiceURL+"/ocr", &requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request to OCR service"})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling OCR service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call OCR service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("OCR service returned non-OK status: %s, body: %s", resp.Status, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OCR service returned an error", "details": string(bodyBytes)})
		return
	}

	// 5. Parse the JSON response from the OCR service
	var ocrResult struct {
		Montant   float64 `json:"montant"`
		Date      string  `json:"date"`
		Tvac      float64 `json:"tvac"`
		Marchand  string  `json:"marchand"`
		TypeRecu  string  `json:"type_recu"`
		IsAnomaly bool    `json:"anomalie"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ocrResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse OCR service response"})
		return
	}

	// 6. Save the data to the database
	parsedDate, err := time.Parse("2006-01-02", ocrResult.Date)
	if err != nil {
		// Default to now if date parsing fails
		parsedDate = time.Now()
	}

	newReceipt := models.Receipt{
		UserID:      currentUser.ID,
		Amount:      ocrResult.Montant,
		Date:        parsedDate,
		Vat:         ocrResult.Tvac,
		Merchant:    ocrResult.Marchand,
		ReceiptType: ocrResult.TypeRecu,
		IsAnomaly:   ocrResult.IsAnomaly,
		FileURL:     "", // We no longer store the file locally
		RawOcrData:  "Data processed by external OCR service",
	}

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
