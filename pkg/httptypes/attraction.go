package httptypes

// Attraction is the info needed for guests to visit an attraction
type Attraction struct {
	URL  string  `json:"url"`
	Fee  float64 `json:"fee"`
	Size float64 `json:"size"` // Size in acres
}
