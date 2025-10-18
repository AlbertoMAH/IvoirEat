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
)

// ReceiptData defines the structure for the data we want the AI to extract.
// We use JSON tags to guide the model's output.
type ReceiptData struct {
	Merchant    string  `json:"merchant"`
	Total       float64 `json:"total"`
	Date        string  `json:"date"` // Expecting "YYYY-MM-DD"
	Category    string  `json:"category"`
	Description string  `json:"description"`
}

// Define the Genkit flow for extracting receipt data.
// This flow is defined once and can be invoked from our HTTP handler.
var extractReceiptFlow = genkit.DefineFlow(
	"extractReceiptFlow",
	func(ctx context.Context, imageBytes []byte, mimeType string) (ReceiptData, error) {
		// Define the prompt for the multimodal model.
		// It includes instructions, the expected JSON format, and the image data.
		prompt := []*ai.Part{
			ai.NewTextPart("Extract the following information from the receipt image and return it as a valid JSON object. " +
				"The fields are: merchant (string), total (float64), date (string, as YYYY-MM-DD), " +
				"category (string, one of: Transport, Restaurant, Hébergement, Achats, Services, Santé, Divers), " +
				"and description (string, a brief summary of the items)."),
			ai.NewDataPart(imageBytes, mimeType),
		}

		// Configure the request to the generative model.
		// We specify the model and ask for JSON output.
		g := genkit.FromContext(ctx)
		resp, err := genkit.Generate(ctx, g,
			ai.WithModelName("googleai/gemini-2.5-flash"),
			ai.WithPrompt(prompt),
			ai.WithOutputFormat(ai.OutputFormatJSON),
		)
		if err != nil {
			return ReceiptData{}, err
		}

		// Extract the JSON text and unmarshal it into our struct.
		var extractedData ReceiptData
		if jsonText, err := resp.Text(); err == nil {
			if err := json.Unmarshal([]byte(jsonText), &extractedData); err != nil {
				return ReceiptData{}, err
			}
		} else {
			return ReceiptData{}, err
		}

		return extractedData, nil
	},
)

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

	// 2. Run the Genkit flow to extract data
	extractedData, err := extractReceiptFlow.Run(c.Request.Context(), fileContents, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to extract data from receipt using AI",
			"details": err.Error(),
		})
		return
	}

	// 3. Parse the date and create the database model
	parsedDate, err := time.Parse("2006-01-02", extractedData.Date)
	if err != nil {
		// If AI gives a bad date format, default to now, but log the issue.
		// In a real app, you might want better error handling.
		parsedDate = time.Now()
	}

	// 4. Save the new receipt to the database
	newReceipt := models.Receipt{
		UserID:      currentUser.ID,
		Amount:      extractedData.Total,
		Date:        parsedDate,
		Vat:         0.0, // VAT is not extracted by this flow for now
		Merchant:    extractedData.Merchant,
		ReceiptType: extractedData.Category,
		Description: extractedData.Description,
		IsAnomaly:   false,
		FileURL:     "", // File storage URL would be set here
		// We can store the AI output as raw data for auditing/debugging
		RawOcrData: string(""), // No longer relevant, could store the JSON output instead
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

	// First, find the receipt to ensure it belongs to the current user
	var receipt models.Receipt
	if err := DB.Where("id = ? AND user_id = ?", receiptID, currentUser.ID).First(&receipt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	// If found, delete it
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

	// Bind the JSON input to a temporary struct
	var input models.Receipt
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the receipt fields
	if err := DB.Model(&receipt).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receipt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"receipt": receipt})
}
