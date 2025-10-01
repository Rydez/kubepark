package main

import (
	"log/slog"
	"time"

	"kubepark/attractions/base"
)

// Carousel represents a carousel attraction
type Carousel struct {
	*base.Attraction
}

// New creates a new carousel attraction
func New() *Carousel {
	config := &base.Config{
		Name:       "carousel",
		Duration:   3 * time.Second,
		BuildCost:  20000,
		RepairCost: 1000,
		Size:       10,
	}

	defaultFee := 5.0

	attraction := base.New(config, defaultFee, afterUse)
	return &Carousel{Attraction: attraction}
}

// afterUse is called after using the carousel
func afterUse() error {
	slog.Debug("Cleaning up carousel after use")
	return nil
}

func main() {
	carousel := New()
	slog.Info("Starting carousel attraction", "park_url", carousel.Config.ParkURL)
	if err := carousel.Start(); err != nil {
		slog.Error("Carousel attraction failed to start", "error", err)
		panic(err)
	}
}
