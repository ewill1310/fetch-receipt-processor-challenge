package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// In-Memory Storage
var receipts = make(map[string]Receipt)
var points = make(map[string]int)

// Process Receipt Endpoint - POST /receipts/process
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	receipts[id] = receipt
	points[id] = calculatePoints(receipt)

	response := ReceiptResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get Points Endpoint - GET /receipts/{id}/points
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	if pts, exists := points[id]; exists {
		response := PointsResponse{Points: pts}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Receipt not found", http.StatusNotFound)
	}
}

// Points Calculation Function
func calculatePoints(receipt Receipt) int {
	totalPoints := 0

	// 1 point per alphanumeric character in retailer name
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	totalPoints += len(re.FindAllString(receipt.Retailer, -1))

	// 50 points if total is a round dollar amount
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == math.Floor(total) {
		totalPoints += 50
	}

	// 25 points if total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		totalPoints += 25
	}

	// 5 points for every two items
	totalPoints += (len(receipt.Items) / 2) * 5

	// Extra points if item description length is a multiple of 3
	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			totalPoints += int(math.Ceil(price * 0.2))
		}
	}

	// 6 points if purchase date is an odd day
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 == 1 {
		totalPoints += 6
	}

	// 10 points if purchase time is between 2:00pm and 4:00pm
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() == 14 {
		totalPoints += 10
	}

	return totalPoints
}

// Start the HTTP Server
func main() {
	http.HandleFunc("/receipts/process", processReceiptHandler)
	http.HandleFunc("/receipts/", getPointsHandler)

	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
