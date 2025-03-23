package main

import (
	"encoding/json"
	"net/http"

	"kubepark/pkg/models"
	"kubepark/pkg/state"
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
func handlePayPark(state *state.GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.PayParkRequest
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
func handleEnterPark(state *state.GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check if park is closed
		if state.IsClosed() {
			json.NewEncoder(w).Encode(models.EnterParkResponse{
				Success: false,
				Message: "Park is closed",
			})
			return
		}

		// Check if park is at capacity
		if !state.CanAddGuest() {
			json.NewEncoder(w).Encode(models.EnterParkResponse{
				Success: false,
				Message: "Park is at capacity",
			})
			return
		}

		// Process entrance fee
		state.AddMoney(state.EntranceFee)

		json.NewEncoder(w).Encode(models.EnterParkResponse{
			Success: true,
			Message: "Welcome to the park!",
		})
	}
}

// handleListAttractions returns a list of available attractions
func handleListAttractions(state *state.GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		attractionStates := state.GetAttractions()
		attractions := make([]models.AttractionInfo, 0, len(attractionStates))
		for _, attractionState := range attractionStates {
			if attractionState.IsRepaired && !attractionState.IsPending {
				attractions = append(attractions, models.AttractionInfo{
					URL: attractionState.URL,
				})
			}
		}

		json.NewEncoder(w).Encode(attractions)
	}
}

// handleRegisterAttraction handles attraction registration requests
func handleRegisterAttraction(state *state.GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.RegisterAttractionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		success, err := state.RegisterAttraction(req)
		if err != nil {
			json.NewEncoder(w).Encode(models.RegisterAttractionResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(models.RegisterAttractionResponse{
			Success: success,
			Message: "Attraction registered successfully",
		})
	}
}
