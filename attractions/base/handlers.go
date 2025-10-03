package base

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"kubepark/pkg/httptypes"
)

// handleAttractionStatus handles the attraction-status endpoint
func handleAttractionStatus(config *Config, state *StateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Return the attraction's fee
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httptypes.Attraction{
			Fee:  config.Fee,
			Size: config.Size,
		})
	}
}

// HandleUse handles the common use endpoint functionality
func handleUse(config *Config, state *StateManager, afterUse func() error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if state.IsBroken() {
			Metrics.AttractionAttempts.WithLabelValues("false", "attraction_broken").Inc()
			http.Error(w, fmt.Sprintf("%s is broken", config.Name), http.StatusServiceUnavailable)
			return
		}

		if config.Closed {
			Metrics.AttractionAttempts.WithLabelValues("false", "attraction_closed").Inc()
			http.Error(w, fmt.Sprintf("%s is closed", config.Name), http.StatusServiceUnavailable)
			return
		}

		// Process payment with kubepark
		if err := ParkTransaction(config, config.Fee); err != nil {
			slog.Error("Failed to process payment", "error", err)
			Metrics.AttractionAttempts.WithLabelValues("false", "payment_failed").Inc()
			http.Error(w, "Payment failed", http.StatusInternalServerError)
			return
		}

		// Simulate usage duration
		time.Sleep(config.Duration)

		// Call after use hook if set
		if afterUse != nil {
			if err := afterUse(); err != nil {
				slog.Error("After use hook failed", "error", err)
				Metrics.AttractionAttempts.WithLabelValues("false", "hook_failed").Inc()
				http.Error(w, "Failed to cleanup attraction", http.StatusInternalServerError)
				return
			}
		}

		Metrics.AttractionAttempts.WithLabelValues("true", "success").Inc()
		w.WriteHeader(http.StatusOK)
	}
}
