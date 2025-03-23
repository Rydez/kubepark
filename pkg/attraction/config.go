package attraction

import (
	"flag"
	"time"
)

// Config represents the common configuration for all attractions
type Config struct {
	Closed     bool
	Fee        float64
	ParkURL    string
	SelfURL    string
	Name       string
	Duration   time.Duration
	BuildCost  int
	RepairCost int
	Size       float64 // Size in acres
}

func RegisterFlags(config *Config, defaultFee float64) {
	flag.BoolVar(&config.Closed, "closed", false, "Whether the attraction is closed")
	flag.StringVar(&config.ParkURL, "park-url", "http://kubepark:80", "URL of the kubepark service")
	flag.StringVar(&config.SelfURL, "self-url", "", "URL where this attraction can be reached")
	flag.Float64Var(&config.Fee, "fee", defaultFee, "Fee for using the attraction")
	flag.Parse()
}
