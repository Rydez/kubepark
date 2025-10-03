package httptypes

type Park struct {
	IsClosed   bool    `json:"is_closed"`   // Whether the park is closed
	TotalSpace float64 `json:"total_space"` // Total space in acres
	Money      float64 `json:"money"`       // Money in the park
}

// TransactionRequest represents a request to send a payment to the park
type TransactionRequest struct {
	Amount float64 `json:"amount"`
}
