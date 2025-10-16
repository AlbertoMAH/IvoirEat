package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/gin-gonic/gin"
	"gobackend/models"
)

// --- Logic ported from Python ---

func findTotal(text string) float64 {
    patterns := []string{
        `(?i)(?:TOTAL|TOTAL\s*TTC|MONTANT\s*TTC|TOTAL\s*A\s*PAYER)\s*[:\s]*([\d\s,.]+)`,
        `(?i)([\d\s,.]+)\s*(?:TOTAL|TTC)`,
    }
    highestAmount := 0.0

    for _, pattern := range patterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindAllStringSubmatch(text, -1)
        for _, match := range matches {
            if len(match) > 1 {
                amountStr := strings.Replace(match[1], ",", ".", -1)
                amountStr = regexp.MustCompile(`[^\d.]`).ReplaceAllString(amountStr, "")
                if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
                    if amount > highestAmount {
                        highestAmount = amount
                    }
                }
            }
        }
    }

    if highestAmount == 0.0 {
        re := regexp.MustCompile(`(\d+[,.]\d{2})`)
        matches := re.FindAllString(text, -1)
        for _, amountStr := range matches {
            amountStr = strings.Replace(amountStr, ",", ".", -1)
            if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
                if amount > highestAmount {
                    highestAmount = amount
                }
            }
        }
    }
    return highestAmount
}

func findDate(text string) string {
    re := regexp.MustCompile(`(\d{2}[-/]\d{2}[-/]\d{4})|(\d{4}[-/]\d{2}[-/]\d{2})`)
    match := re.FindString(text)
    if match != "" {
        // Normalize date to YYYY-MM-DD
        layouts := []string{"02-01-2006", "02/01/2006", "2006-01-02"}
        for _, layout := range layouts {
            if t, err := time.Parse(layout, match); err == nil {
                return t.Format("2006-01-02")
            }
        }
    }
    return time.Now().Format("2006-01-02")
}

func findMerchant(text string) string {
    lines := strings.Split(text, "\n")
    if len(lines) > 0 {
        for _, line := range lines {
            trimmedLine := strings.TrimSpace(line)
            if trimmedLine != "" {
                return trimmedLine
            }
        }
    }
    return "Unknown Merchant"
}

func classifyReceipt(text string) string {
    textLower := strings.ToLower(text)
    categories := map[string][]string{
        "Transport":  {"taxi", "trajet", "course", "station", "carburant", "ticket", "vol"},
        "Restaurant": {"restaurant", "repas", "menu", "table", "boisson", "cafe", "snack"},
        "Hébergement": {"hotel", "chambre", "nuitée", "séjour"},
        "Achats":     {"supermarché", "articles", "ttc", "boutique", "fournitures"},
        "Services":   {"abonnement", "licence", "honoraires", "consultation", "réparation"},
        "Santé":      {"pharmacie", "medicament", "consultation", "soin"},
    }

    for category, keywords := range categories {
        for _, keyword := range keywords {
            if strings.Contains(textLower, keyword) {
                return category
            }
        }
    }
    return "Divers"
}

// --- Main Controller Function ---

func UploadReceipt(c *gin.Context) {
	// 1. Get user and file
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

	// 2. Call Google Cloud Vision API
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Vision client", "details": err.Error()})
		return
	}
	defer client.Close()

	image, err := vision.NewImageFromReader(strings.NewReader(string(fileContents)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image from reader", "details": err.Error()})
		return
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect text", "details": err.Error()})
		return
	}

	var rawText string
	if len(annotations) > 0 {
		rawText = annotations[0].Description
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No text found in image"})
		return
	}

	// 3. Extract data from text
	totalAmount := findTotal(rawText)
	receiptDate := findDate(rawText)
	merchantName := findMerchant(rawText)
	receiptType := classifyReceipt(rawText)
	parsedDate, _ := time.Parse("2006-01-02", receiptDate)

	// 4. Save to database
	newReceipt := models.Receipt{
		UserID:      currentUser.ID,
		Amount:      totalAmount,
		Date:        parsedDate,
		Vat:         0.0,
		Merchant:    merchantName,
		ReceiptType: receiptType,
		IsAnomaly:   false,
		FileURL:     "",
		RawOcrData:  rawText,
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
	DB.Where("user_id = ?", currentUser.ID).Find(&receipts)
	c.JSON(http.StatusOK, gin.H{"receipts": receipts})
}
