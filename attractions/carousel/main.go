package main

import (
	"log"
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
	}

	defaultFee := 5.0

	attraction := base.New(config, defaultFee, afterUse)
	return &Carousel{Attraction: attraction}
}

// afterUse is called after using the carousel
func afterUse() error {
	log.Printf("Cleaning up carousel after use")
	return nil
}

func main() {
	carousel := New()
	log.Printf("Starting carousel attraction with park URL: %s", carousel.Config.ParkURL)
	log.Fatal(carousel.Start())
}
