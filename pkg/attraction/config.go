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
	Name       string
	Duration   time.Duration
	BuildCost  int
	RepairCost int
}

func RegisterFlags(config *Config, defaultFee float64) {
	flag.BoolVar(&config.Closed, "closed", false, "Whether the attraction is closed")
	flag.StringVar(&config.ParkURL, "park-url", "http://kubepark:80", "URL of the kubepark service")
	flag.Float64Var(&config.Fee, "fee", defaultFee, "Fee for using the attraction")
	flag.Parse()
}
