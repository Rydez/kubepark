package main

import (
	"log"
	"net/http"
	"time"

	"kubepark/pkg/attraction"
	"kubepark/pkg/metrics"
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

	if r.IsClosed() {
		metrics.AttractionAttempts.WithLabelValues("false", "attraction_closed").Inc()
		http.Error(w, "Restroom is closed", http.StatusServiceUnavailable)
		return
	}

	// Process payment with kubepark
	if err := r.ProcessPayment(); err != nil {
		log.Printf("Failed to process payment: %v", err)
		http.Error(w, "Payment failed", http.StatusInternalServerError)
		return
	}

	// Simulate usage duration
	time.Sleep(r.GetDuration())

	metrics.Revenue.Add(r.GetFee())
	metrics.AttractionAttempts.WithLabelValues("true", "success").Inc()

	w.WriteHeader(http.StatusOK)
}

func main() {
	restroom := New()
	log.Printf("Starting restroom attraction with park URL: %s", restroom.GetParkURL())
	log.Fatal(restroom.Start())
}
