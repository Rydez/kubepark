package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"kubepark/pkg/httptypes"
)

// handleStatus handles requests to check if this is a park service
func handleStatus(config *Config, state *StateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httptypes.Park{
			IsClosed:   isClosed(config, state.GetTime()),
			TotalSpace: state.GetTotalSpace(),
			Money:      state.GetMoney(),
		})
	}
}

// handleTransaction handles payment requests from attractions
func handleTransaction(state *StateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req httptypes.TransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Add the payment to the park's cash
		state.AddMoney(req.Amount)
		w.WriteHeader(http.StatusOK)
	}
}

// handleEnter handles guest entry requests
func handleEnter(state *StateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Process entrance fee
		state.AddMoney(state.GetEntranceFee())
		slog.Info("Accepted guest")
		w.WriteHeader(http.StatusOK)
	}
}
