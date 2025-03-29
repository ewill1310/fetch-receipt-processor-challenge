# Receipt Processor Web Service (Go)
The service was initially coded/ran on a windows device, with no current access to MacOS or Linux so testing for those platforms have not been done. 
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
$headers = @{
    "Content-Type" = "application/json"
}

$response = Invoke-WebRequest -Uri "http://localhost:8080/receipts/process" -Method Post -Headers $headers -Body '{
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
        { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
        { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
    ],
    "total": "35.35"
}'

$response.Content
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
$response = Invoke-WebRequest -Uri "http://localhost:8080/receipts/6f985bef-4c8a-4fae-8dfb-ecdb0a4240ce/points" -Method Get

$response.Content
```

Example Response
```sh
{ "points": 28 }
```
### 3. Stopping the Server
To stop the server, simply press Ctrl + C in the terminal where it's running.

## Testing
You can manually test the API using powershell, curl, Postman, or any API testing tool.

## Author:
Evan Williams
