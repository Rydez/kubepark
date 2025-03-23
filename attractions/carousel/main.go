package main

import (
	"log"
	"net/http"
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

	attraction := base.New(config, defaultFee)
	carousel := &Carousel{Attraction: attraction}

	// Add the /use endpoint specific to the carousel
	mainMux := carousel.MainServer.Handler.(*http.ServeMux)
	mainMux.HandleFunc("/use", carousel.handleUse)

	return carousel
}

// handleUse handles the carousel's use endpoint
func (c *Carousel) handleUse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if c.Config.Closed {
		base.Metrics.AttractionAttempts.WithLabelValues("false", "attraction_closed").Inc()
		http.Error(w, "Carousel is closed", http.StatusServiceUnavailable)
		return
	}

	// Validate payment
	_, err := c.ValidatePayment(w, r)
	if err != nil {
		base.Metrics.AttractionAttempts.WithLabelValues("false", "insufficient_funds").Inc()
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}

	// Process payment with kubepark
	if err := c.PayPark(); err != nil {
		log.Printf("Failed to process payment: %v", err)
		base.Metrics.AttractionAttempts.WithLabelValues("false", "payment_failed").Inc()
		http.Error(w, "Payment failed", http.StatusInternalServerError)
		return
	}

	// Simulate ride duration
	time.Sleep(c.Config.Duration)

	base.Metrics.AttractionAttempts.WithLabelValues("true", "success").Inc()
	w.WriteHeader(http.StatusOK)
}

func main() {
	carousel := New()
	log.Printf("Starting carousel attraction with park URL: %s", carousel.Config.ParkURL)
	log.Fatal(carousel.Start())
}
