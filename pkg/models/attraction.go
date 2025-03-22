package models

// AttractionInfo represents the public information about an attraction
type AttractionInfo struct {
	URL string `json:"url"`
}

// RegisterAttractionRequest represents a request to register a new attraction
type RegisterAttractionRequest struct {
	URL        string `json:"url"`
	BuildCost  int    `json:"build_cost"`
	RepairCost int    `json:"repair_cost"`
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
