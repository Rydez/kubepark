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

// BreakAttractionRequest represents a request to break the attraction
type BreakAttractionRequest struct {
	URL string `json:"url"`
}

// IsAttractionResponse represents the response from the /is-attraction endpoint
type IsAttractionResponse struct {
	Fee float64 `json:"fee"`
}
