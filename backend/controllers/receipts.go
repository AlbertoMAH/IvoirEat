package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gobackend/models"
)

// --- Local OCR Service Response Structure ---
type LocalOCRResponse struct {
	Filename string `json:"filename"`
	Text     string `json:"text"`
}

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
	// 1. Get user and file from the request
	user, _ := c.Get("user")
	currentUser := user.(models.User)
	fileHeader, err := c.FormFile("receiptImage")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed"})
		return
	}

	// 2. Prepare the request to our local OCR service
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	file, _ := fileHeader.Open()
	part, _ := writer.CreateFormFile("file", fileHeader.Filename)
	_, _ = io.Copy(part, file)
	file.Close()
	writer.Close()

	// 3. Send the request to the OCR service container
	// We use the service name 'ocr_service' as the host, which is resolved by Docker's internal DNS
	ocrServiceURL := "http://ocr_service:8000/ocr"
	req, _ := http.NewRequest("POST", ocrServiceURL, &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set a generous timeout as EasyOCR can be slow, especially on first run
	client := &http.Client{Timeout: time.Minute * 2}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call local OCR service", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OCR service returned an error", "details": string(bodyBytes)})
		return
	}

	// 4. Parse the response from our local OCR service
	var ocrResponse LocalOCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse OCR response"})
		return
	}
	rawText := ocrResponse.Text

	// 5. Extract data from text
	totalAmount := findTotal(rawText)
	receiptDate := findDate(rawText)
	merchantName := findMerchant(rawText)
	receiptType := classifyReceipt(rawText)
	parsedDate, _ := time.Parse("2006-01-02", receiptDate)

	// 6. Save the new receipt to the database
	newReceipt := models.Receipt{
		UserID:      currentUser.ID,
		Amount:      totalAmount,
		Date:        parsedDate,
		Vat:         0.0, // VAT logic not implemented yet
		Merchant:    merchantName,
		ReceiptType: receiptType,
		IsAnomaly:   false, // Anomaly logic not implemented yet
		FileURL:     "",    // No longer storing file
		RawOcrData:  rawText,
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
	user, _ := c.Get("user")
	currentUser := user.(models.User)
	var receipts []models.Receipt
	DB.Where("user_id = ?", currentUser.ID).Find(&receipts)
	c.JSON(http.StatusOK, gin.H{"receipts": receipts})
}
