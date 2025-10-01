package main

import (
	"log/slog"
	"time"

	"kubepark/attractions/base"
)

// Restroom represents a restroom attraction
type Restroom struct {
	*base.Attraction
}

// New creates a new restroom attraction
func New() *Restroom {
	config := &base.Config{
		Name:       "restroom",
		Duration:   2 * time.Second,
		BuildCost:  10000,
		RepairCost: 500,
		Size:       1,
	}

	defaultFee := 2.0

	attraction := base.New(config, defaultFee, afterUse)
	return &Restroom{Attraction: attraction}
}

// afterUse is called after using the restroom
func afterUse() error {
	slog.Debug("Cleaning up restroom after use")
	return nil
}

func main() {
	restroom := New()
	slog.Info("Starting restroom attraction", "park_url", restroom.Config.ParkURL)
	if err := restroom.Start(); err != nil {
		slog.Error("Restroom attraction failed to start", "error", err)
		panic(err)
	}
}
