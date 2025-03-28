package httptypes

// Attraction is the info needed for guests to visit an attraction
type Attraction struct {
	URL string `json:"url"`
}

// RegisterAttractionRequest represents a request to register a new attraction
type RegisterAttractionRequest struct {
	URL        string  `json:"url"`
	BuildCost  int     `json:"build_cost"`
	RepairCost int     `json:"repair_cost"`
	Size       float64 `json:"size"` // Size in acres
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

// BreakAttractionRequest represents a request to break the attraction
type BreakAttractionRequest struct {
	URL string `json:"url"`
}

// BreakAttractionResponse represents the response from the break endpoint
type BreakAttractionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
