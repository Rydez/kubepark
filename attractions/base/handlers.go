package base

import (
	"encoding/json"
	"net/http"

	"kubepark/pkg/httptypes"
)

// handleIsAttraction handles the is-attraction endpoint
func handleIsAttraction(config *Config, state *StateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Return the attraction's fee
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httptypes.IsAttractionResponse{
			Fee: config.Fee,
		})
	}
}
