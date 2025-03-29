package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Item represents an item in a receipt.
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Receipt represents a full receipt with purchase details.
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// ReceiptResponse represents the response when a receipt is processed.
type ReceiptResponse struct {
	ID string `json:"id"`
}

// PointsResponse represents the response when retrieving receipt points.
type PointsResponse struct {
	Points int `json:"points"`
}

// In-memory storage for receipts.
var (
	receipts = make(map[string]Receipt)
	mu       sync.Mutex // Ensures thread safety
)

// processReceipt handles the POST request to process a receipt.
func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the receipt
	id := uuid.New().String()

	// Store the receipt in memory
	mu.Lock()
	receipts[id] = receipt
	mu.Unlock()

	// Send response
	response := ReceiptResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getPoints handles the GET request to retrieve points for a receipt.
func getPoints(w http.ResponseWriter, r *http.Request) {
	// Extract receipt ID from the URL
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	// Retrieve receipt from memory
	mu.Lock()
	receipt, exists := receipts[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Calculate points
	points := calculatePoints(receipt)

	// Send response
	response := PointsResponse{Points: points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// calculatePoints applies the scoring rules to a receipt.
func calculatePoints(receipt Receipt) int {
	points := 0

	// 1. One point for every alphanumeric character in the retailer name
	alphanumericRegex := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(alphanumericRegex.FindAllString(receipt.Retailer, -1))

	// 2. 50 points if the total is a round dollar amount with no cents
	if total, err := strconv.ParseFloat(receipt.Total, 64); err == nil && total == math.Floor(total) {
		points += 50
	}

	// 3. 25 points if the total is a multiple of 0.25
	if total, err := strconv.ParseFloat(receipt.Total, 64); err == nil && math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// 4. 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// 5. Points based on item description length
	for _, item := range receipt.Items {
		description := strings.TrimSpace(item.ShortDescription)
		if len(description)%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// 6. 6 points if the day in the purchase date is odd
	if dateParts := strings.Split(receipt.PurchaseDate, "-"); len(dateParts) == 3 {
		if day, err := strconv.Atoi(dateParts[2]); err == nil && day%2 == 1 {
			points += 6
		}
	}

	// 7. 10 points if the purchase time is between 2:00pm and 4:00pm
	if timeParts := strings.Split(receipt.PurchaseTime, ":"); len(timeParts) == 2 {
		if hour, err := strconv.Atoi(timeParts[0]); err == nil {
			if minute, err := strconv.Atoi(timeParts[1]); err == nil {
				purchaseTime := hour*60 + minute
				if purchaseTime >= 840 && purchaseTime < 960 { // 2:00pm = 840 min, 4:00pm = 960 min
					points += 10
				}
			}
		}
	}

	return points
}

func main() {
	http.HandleFunc("/receipts/process", processReceipt)
	http.HandleFunc("/receipts/", getPoints)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
