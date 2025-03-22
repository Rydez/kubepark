package main

import (
	"log"
	"net/http"
	"time"

	"kubepark/pkg/attraction"
	"kubepark/pkg/metrics"
)

// Carousel represents a carousel attraction
type Carousel struct {
	*attraction.Attraction
}

// New creates a new carousel attraction
func New() *Carousel {
	base := attraction.New("carousel", 5, 3*time.Second)
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

	if c.IsClosed() {
		metrics.AttractionAttempts.WithLabelValues("false", "attraction_closed").Inc()
		http.Error(w, "Carousel is closed", http.StatusServiceUnavailable)
		return
	}

	// Process payment with kubepark
	if err := c.ProcessPayment(); err != nil {
		log.Printf("Failed to process payment: %v", err)
		http.Error(w, "Payment failed", http.StatusInternalServerError)
		return
	}

	// Simulate ride duration
	time.Sleep(c.GetDuration())

	metrics.Revenue.Add(c.GetFee())
	metrics.AttractionAttempts.WithLabelValues("true", "success").Inc()

	w.WriteHeader(http.StatusOK)
}

func main() {
	carousel := New()
	log.Printf("Starting carousel attraction with park URL: %s", carousel.GetParkURL())
	log.Fatal(carousel.Start())
}
