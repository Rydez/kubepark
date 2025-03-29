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
func handlePayPark(state *GameState) http.HandlerFunc {
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
func handleEnterPark(state *GameState) http.HandlerFunc {
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
		state.AddMoney(state.EntranceFee)
		log.Printf("Accepted guest")
		w.WriteHeader(http.StatusOK)
	}
}

// handleListAttractions returns a list of available attractions
func handleListAttractions(state *GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		attractions := []httptypes.Attraction{}
		attractionStates := state.GetAttractions()
		for _, attractionState := range attractionStates {
			if attractionState.IsRepaired && !attractionState.IsPending {
				attractions = append(attractions, httptypes.Attraction{
					URL: attractionState.URL,
				})
			}
		}

		json.NewEncoder(w).Encode(attractions)
	}
}

// handleRegisterAttraction handles attraction registration requests
func handleRegisterAttraction(state *GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req httptypes.RegisterAttractionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := state.RegisterAttraction(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleBreakAttraction(state *GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req httptypes.BreakAttractionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := state.MarkAttractionBroken(req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
