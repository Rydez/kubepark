package main

import (
	"log"
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
	log.Printf("Cleaning up restroom after use")
	return nil
}

func main() {
	restroom := New()
	log.Printf("Starting restroom attraction with park URL: %s", restroom.Config.ParkURL)
	log.Fatal(restroom.Start())
}
