package httptypes

// PayParkRequest represents a request to send a payment to the park
type PayParkRequest struct {
	Amount float64 `json:"amount"`
}
