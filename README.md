# Receipt Processor Web Service (Go)

## Overview
This is a RESTful web service that processes receipts and calculates reward points based on predefined rules. The service provides two endpoints:

1. `POST /receipts/process` - Accepts a JSON receipt, generates a unique ID, and calculates reward points.
2. `GET /receipts/{id}/points` - Returns the reward points for a given receipt ID.

## Technologies Used
- **Language**: Go (Golang)
- **Framework**: Standard `net/http` package
- **Data Storage**: In-memory (no database required)
- **UUID Generator**: `github.com/google/uuid`

## Installation & Running the Service

### **1. Clone the Repository**
```sh
git clone https://github.com/YOUR_GITHUB_USERNAME/receipt-processor.git
cd receipt-processor
```
### 2. Install Dependencies
```sh
go mod tidy
```
### 3. Run the Service
```sh
go run main.go
```
The server will start on http://localhost:8080.

## API Endpoints
### 1. Process a Receipt
Endpoint: POST /receipts/process

Description: Accepts a receipt JSON, assigns it a unique ID, and calculates points.

Example Request
```sh
curl -X POST http://localhost:8080/receipts/process \
     -H "Content-Type: application/json" \
     -d '{
           "retailer": "Target",
           "purchaseDate": "2022-01-01",
           "purchaseTime": "13:01",
           "items": [
             { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
             { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
           ],
           "total": "18.74"
         }'
```

Example Response
```sh
{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
```

### 2. Get Receipt Points
Endpoint: GET /receipts/{id}/points

Description: Retrieves the points for a given receipt ID.

Example Request
```sh
curl -X GET http://localhost:8080/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points
```

Example Response
```sh
{ "points": 28 }
```

## Testing
You can manually test the API using curl, Postman, or any API testing tool.

## Author:
Evan Williams
