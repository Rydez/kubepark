package attraction

import (
	"flag"
	"time"
)

// Config represents the common configuration for all attractions
type Config struct {
	Closed     bool
	Fee        float64
	VolumePath string
	ParkURL    string
	Name       string
	Duration   time.Duration
}

func RegisterFlags(config *Config, defaultFee float64) {
	flag.BoolVar(&config.Closed, "closed", false, "Whether the attraction is closed")
	flag.StringVar(&config.ParkURL, "park-url", "http://kubepark:80", "URL of the kubepark service")
	flag.StringVar(&config.VolumePath, "volume", "", "Path to volume for persistent storage")
	flag.Float64Var(&config.Fee, "fee", defaultFee, "Fee for using the attraction")
	flag.Parse()
}
