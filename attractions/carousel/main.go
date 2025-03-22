package main

import (
	"log"
	"net/http"
	"time"

	"kubepark/pkg/attraction"
)

// Carousel represents a carousel attraction
type Carousel struct {
	*attraction.Attraction
}

// New creates a new carousel attraction
func New() *Carousel {
	config := &attraction.Config{
		Name:       "carousel",
		Duration:   3 * time.Second,
		BuildCost:  20000,
		RepairCost: 1000,
	}

	defaultFee := 5.0

	base := attraction.New(config, defaultFee)
	carousel := &Carousel{Attraction: base}

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
		attraction.Metrics.AttractionAttempts.WithLabelValues("false", "attraction_closed").Inc()
		http.Error(w, "Carousel is closed", http.StatusServiceUnavailable)
		return
	}

	// Validate payment
	_, err := c.ValidatePayment(w, r)
	if err != nil {
		attraction.Metrics.AttractionAttempts.WithLabelValues("false", "insufficient_funds").Inc()
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}

	// Process payment with kubepark
	if err := c.PayPark(); err != nil {
		log.Printf("Failed to process payment: %v", err)
		attraction.Metrics.AttractionAttempts.WithLabelValues("false", "payment_failed").Inc()
		http.Error(w, "Payment failed", http.StatusInternalServerError)
		return
	}

	// Simulate ride duration
	time.Sleep(c.Config.Duration)

	attraction.Metrics.Revenue.Add(c.Config.Fee)
	attraction.Metrics.AttractionAttempts.WithLabelValues("true", "success").Inc()

	w.WriteHeader(http.StatusOK)
}

func main() {
	carousel := New()
	log.Printf("Starting carousel attraction with park URL: %s", carousel.Config.ParkURL)
	log.Fatal(carousel.Start())
}
