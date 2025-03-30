package main

import (
	"encoding/json"
	"log"
	"net/http"

	"kubepark/pkg/httptypes"
)

// handleParkStatus handles requests to check if this is a park service
func handleParkStatus(state *StateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httptypes.Park{
			TotalSpace: state.GetTotalSpace(),
			Money:      state.GetMoney(),
		})
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
