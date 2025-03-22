package models

// AttractionInfo represents the public information about an attraction
type AttractionInfo struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	URL         string  `json:"url"`
}

// ListAttractionInfoResponse represents the response from the attractions endpoint
type ListAttractionInfoResponse struct {
	Attractions []AttractionInfo `json:"attractions"`
}

// RegisterAttractionRequest represents a request to register a new attraction
type RegisterAttractionRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	RepairFee   float64 `json:"repair_fee"`
	URL         string  `json:"url"`
}

// RegisterAttractionResponse represents the response from the register endpoint
type RegisterAttractionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UseAttractionRequest represents a request to use the attraction
type UseAttractionRequest struct {
	GuestMoney float64 `json:"guest_money"`
}
