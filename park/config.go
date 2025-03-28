package main

import (
	"flag"
)

// Config represents the park configuration
type Config struct {
	Image       string
	SelfURL     string
	Mode        string
	VolumePath  string
	Closed      bool
	EntranceFee float64
	OpensAt     int
	ClosesAt    int
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
	flag.Parse()
}
