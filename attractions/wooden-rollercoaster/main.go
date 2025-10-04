package main

import (
	"log/slog"
	"time"

	"kubepark/attractions/base"
)

// WoodenRollercoaster represents a wooden rollercoaster attraction
type WoodenRollercoaster struct {
	*base.Attraction
}

// New creates a new wooden rollercoaster attraction
func New() *WoodenRollercoaster {
	config := &base.Config{
		Name:       "wooden-rollercoaster",
		Duration:   45 * time.Second,
		BuildCost:  150000,
		RepairCost: 5000,
		Size:       25, // Large footprint for a rollercoaster
	}

	defaultFee := 15.0 // Higher fee for thrilling attraction

	attraction := base.New(config, defaultFee, afterUse)
	return &WoodenRollercoaster{Attraction: attraction}
}

// afterUse is called after using the wooden rollercoaster
func afterUse() error {
	slog.Debug("Cleaning up wooden rollercoaster after use")
	return nil
}

func main() {
	woodenRollercoaster := New()
	slog.Info("Starting wooden rollercoaster attraction", "park_url", woodenRollercoaster.Config.ParkURL)
	if err := woodenRollercoaster.Start(); err != nil {
		slog.Error("Wooden rollercoaster attraction failed to start", "error", err)
		panic(err)
	}
}
