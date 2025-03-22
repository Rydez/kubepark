package models

// EnterParkResponse represents the response from the enter endpoint
type EnterParkResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PayParkRequest represents a request to send a payment to the park
type PayParkRequest struct {
	Amount float64 `json:"amount"`
}
