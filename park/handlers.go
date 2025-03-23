package main

import (
	"encoding/json"
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
			json.NewEncoder(w).Encode(httptypes.EnterParkResponse{
				Success: false,
				Message: "Park is closed",
			})
			return
		}

		// Check if park is at capacity
		if !state.CanAddGuest() {
			json.NewEncoder(w).Encode(httptypes.EnterParkResponse{
				Success: false,
				Message: "Park is at capacity",
			})
			return
		}

		// Process entrance fee
		state.AddMoney(state.EntranceFee)

		json.NewEncoder(w).Encode(httptypes.EnterParkResponse{
			Success: true,
			Message: "Welcome to the park!",
		})
	}
}

// handleListAttractions returns a list of available attractions
func handleListAttractions(state *GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := httptypes.ListAttractionsResponse{}
		attractionStates := state.GetAttractions()
		for _, attractionState := range attractionStates {
			if attractionState.IsRepaired && !attractionState.IsPending {
				data.Attractions = append(data.Attractions, httptypes.Attraction{
					URL: attractionState.URL,
				})
			}
		}

		json.NewEncoder(w).Encode(data)
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

		success, err := state.RegisterAttraction(req)
		if err != nil {
			json.NewEncoder(w).Encode(httptypes.RegisterAttractionResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(httptypes.RegisterAttractionResponse{
			Success: success,
			Message: "Attraction registered successfully",
		})
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

		json.NewEncoder(w).Encode(httptypes.BreakAttractionResponse{
			Success: true,
			Message: "Attraction broken",
		})
	}
}
