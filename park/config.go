package main

import (
	"flag"
	"os"
)

// Config represents the park configuration
type Config struct {
	Image         string
	SelfURL       string
	Mode          string
	VolumePath    string
	Closed        bool
	EntranceFee   float64
	OpensAt       int
	ClosesAt      int
	LogLevel      string
	GrafanaURL    string
	GrafanaAPIKey string
}

func RegisterFlags(config *Config) {
	flag.StringVar(&config.Image, "image", "", "The image to use for guests, should be same as the park image")
	flag.StringVar(&config.SelfURL, "self-url", "", "URL where this attraction can be reached")
	flag.StringVar(&config.Mode, "mode", "easy", "Game mode (easy, medium, hard)")
	flag.StringVar(&config.VolumePath, "volume", "", "Path to volume for persistent storage")
	flag.BoolVar(&config.Closed, "closed", false, "Whether the park is closed")
	flag.Float64Var(&config.EntranceFee, "entrance-fee", 10, "Entrance fee for the park")
	flag.IntVar(&config.OpensAt, "opens-at", 8, "Hour at which the park opens")
	flag.IntVar(&config.ClosesAt, "closes-at", 20, "Hour at which the park closes")
	flag.StringVar(&config.LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&config.GrafanaURL, "grafana-url", "http://kubepark-grafana:3000", "Grafana server URL for Live streaming")
	flag.StringVar(&config.GrafanaAPIKey, "grafana-api-key", "", "Grafana API key for Live streaming")
	flag.Parse()

	// Override with environment variables if set
	if envAPIKey := os.Getenv("GRAFANA_API_KEY"); envAPIKey != "" {
		config.GrafanaAPIKey = envAPIKey
	}
}
