package main

import (
	"log"
	"net/http"
	"time"

	"kubepark/pkg/attraction"
)

// Restroom represents a restroom attraction
type Restroom struct {
	*attraction.Attraction
}

// New creates a new restroom attraction
func New() *Restroom {
	base := attraction.New("restroom", 2, 2*time.Second)
	restroom := &Restroom{Attraction: base}

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
		attraction.Metrics.AttractionAttempts.WithLabelValues("false", "attraction_closed").Inc()
		http.Error(w, "Restroom is closed", http.StatusServiceUnavailable)
		return
	}

	// Validate payment
	_, err := r.ValidatePayment(w, req)
	if err != nil {
		attraction.Metrics.AttractionAttempts.WithLabelValues("false", "insufficient_funds").Inc()
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}

	// Process payment with kubepark
	if err := r.PayPark(); err != nil {
		log.Printf("Failed to process payment: %v", err)
		attraction.Metrics.AttractionAttempts.WithLabelValues("false", "payment_failed").Inc()
		http.Error(w, "Payment failed", http.StatusInternalServerError)
		return
	}

	// Simulate usage duration
	time.Sleep(r.Config.Duration)

	attraction.Metrics.Revenue.Add(r.Config.Fee)
	attraction.Metrics.AttractionAttempts.WithLabelValues("true", "success").Inc()

	w.WriteHeader(http.StatusOK)
}

func main() {
	restroom := New()
	log.Printf("Starting restroom attraction with park URL: %s", restroom.Config.ParkURL)
	log.Fatal(restroom.Start())
}
