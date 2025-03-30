package httptypes

type Park struct {
	TotalSpace float64 `json:"total_space"` // Total space in acres
	Money      float64 `json:"money"`       // Money in the park
}

// PayParkRequest represents a request to send a payment to the park
type PayParkRequest struct {
	Amount float64 `json:"amount"`
}
