package main

import (
	"log"
	"net/http"
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
		Fee:        2,
		Duration:   2 * time.Second,
		BuildCost:  10000,
		RepairCost: 500,
	}

	defaultFee := 2.0

	attraction := base.New(config, defaultFee)
	restroom := &Restroom{Attraction: attraction}

	// Add the /use endpoint specific to the restroom
	mainMux := restroom.MainServer.Handler.(*http.ServeMux)
	mainMux.HandleFunc("/use", restroom.handleUse)

	return restroom
}

// handleUse handles the restroom's use endpoint
func (r *Restroom) handleUse(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Config.Closed {
		base.Metrics.AttractionAttempts.WithLabelValues("false", "attraction_closed").Inc()
		http.Error(w, "Restroom is closed", http.StatusServiceUnavailable)
		return
	}

	// Validate payment
	_, err := r.ValidatePayment(w, req)
	if err != nil {
		base.Metrics.AttractionAttempts.WithLabelValues("false", "insufficient_funds").Inc()
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}

	// Process payment with kubepark
	if err := r.PayPark(); err != nil {
		log.Printf("Failed to process payment: %v", err)
		base.Metrics.AttractionAttempts.WithLabelValues("false", "payment_failed").Inc()
		http.Error(w, "Payment failed", http.StatusInternalServerError)
		return
	}

	// Simulate usage duration
	time.Sleep(r.Config.Duration)

	base.Metrics.AttractionAttempts.WithLabelValues("true", "success").Inc()
	w.WriteHeader(http.StatusOK)
}

func main() {
	restroom := New()
	log.Printf("Starting restroom attraction with park URL: %s", restroom.Config.ParkURL)
	log.Fatal(restroom.Start())
}
