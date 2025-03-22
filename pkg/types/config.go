package types

// ParkConfig represents the configuration for the main kubepark simulator
type ParkConfig struct {
	Closed      bool    `json:"closed"`
	EntranceFee float64 `json:"entrance_fee"`
	VolumePath  string  `json:"volume_path,omitempty"`
	Mode        string  `json:"mode"`
	OpensAt     int     `json:"opens_at"`
	ClosesAt    int     `json:"closes_at"`
}

// AttractionConfig represents the base configuration for all attractions
type AttractionConfig struct {
	Closed     bool    `json:"closed"`
	Fee        float64 `json:"fee"`
	VolumePath string  `json:"volume_path,omitempty"`
}
