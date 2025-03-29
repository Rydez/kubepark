package main

import (
	"encoding/json"
	"log"
	"net/http"

	"kubepark/pkg/httptypes"
)

// handleIsPark handles requests to check if this is a park service
func handleIsPark() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// handlePayPark handles payment requests from attractions
func handlePayPark(state *StateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req httptypes.PayParkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Add the payment to the park's cash
		state.AddMoney(req.Amount)
		w.WriteHeader(http.StatusOK)
	}
}

// handleEnterPark handles guest entry requests
func handleEnterPark(state *StateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check if park is closed
		if state.IsClosed() {
			log.Printf("Turned away guest, park is closed")
			http.Error(w, "Park is closed", http.StatusBadRequest)
			return
		}

		// Check if park is at capacity
		if !state.CanAddGuest() {
			log.Printf("Turned away guest, park is at capacity")
			http.Error(w, "Park is at capacity", http.StatusBadRequest)
			return
		}

		// Process entrance fee
		state.AddMoney(state.GetEntranceFee())
		log.Printf("Accepted guest")
		w.WriteHeader(http.StatusOK)
	}
}
