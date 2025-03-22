package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"kubepark/pkg/handlers"
	"kubepark/pkg/metrics"
	"kubepark/pkg/state"
)

// Config represents the park configuration
type Config struct {
	Mode        string
	VolumePath  string
	Closed      bool
	EntranceFee float64
	OpensAt     int
	ClosesAt    int
}

// Park represents the amusement park simulator
type Park struct {
	Config        *Config
	MetricsServer *http.Server
	MainServer    *http.Server
	State         *state.GameState
}

// PaymentRequest represents a request to process a payment
type PaymentRequest struct {
	Amount float64 `json:"amount"`
}

// New creates a new park simulator
func New() *Park {
	config := &Config{}
	flag.StringVar(&config.Mode, "mode", "easy", "Game mode (easy, medium, hard)")
	flag.StringVar(&config.VolumePath, "volume", "", "Path to volume for persistent storage")
	flag.BoolVar(&config.Closed, "closed", false, "Whether the park is closed")
	flag.Float64Var(&config.EntranceFee, "entrance-fee", 10, "Entrance fee for the park")
	flag.IntVar(&config.OpensAt, "opens-at", 8, "Hour at which the park opens")
	flag.IntVar(&config.ClosesAt, "closes-at", 20, "Hour at which the park closes")
	flag.Parse()

	// Initialize game state
	gameState, err := state.New(config.VolumePath)
	if err != nil {
		log.Fatalf("Failed to initialize game state: %v", err)
	}

	// Update state with config values
	gameState.Mode = config.Mode
	gameState.EntranceFee = config.EntranceFee
	gameState.OpensAt = config.OpensAt
	gameState.ClosesAt = config.ClosesAt
	gameState.SetClosed(config.Closed)

	metrics.RegisterParkMetrics()
	metrics.EntranceFee.Set(config.EntranceFee)
	metrics.IsParkClosed.Set(btof(config.Closed))
	metrics.OpensAt.Set(float64(config.OpensAt))
	metrics.ClosesAt.Set(float64(config.ClosesAt))

	// Create metrics server on port 9000
	metricsServer := &http.Server{
		Addr:    ":9000",
		Handler: handlers.NewParkHandler(),
	}

	// Create main server on port 80
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/pay", handlePayment(gameState))
	mainServer := &http.Server{
		Addr:    ":80",
		Handler: mainMux,
	}

	return &Park{
		Config:        config,
		MetricsServer: metricsServer,
		MainServer:    mainServer,
		State:         gameState,
	}
}

// handlePayment handles payment requests from attractions
func handlePayment(state *state.GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req PaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Add the payment to the park's cash
		state.AddCash(req.Amount)
		w.WriteHeader(http.StatusOK)
	}
}

// Start starts the park simulator
func (p *Park) Start() error {
	// Start the metrics server
	go func() {
		if err := p.MetricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Start the park simulation loop
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Update metrics
			metrics.Cash.Set(p.State.GetCash())
			metrics.Time.Set(float64(p.State.GetTime().Unix()))

			// Save state every minute
			if time.Since(p.State.GetTime()) > time.Minute {
				if err := p.State.Save(); err != nil {
					log.Printf("Failed to save state: %v", err)
				}
			}

			// TODO: Create and manage guests
			// TODO: Handle park events
		}
	}()

	// Start the main server
	return p.MainServer.ListenAndServe()
}

// Stop gracefully stops the park simulator
func (p *Park) Stop() error {
	// Save final state
	if err := p.State.Save(); err != nil {
		log.Printf("Failed to save final state: %v", err)
	}

	if err := p.MetricsServer.Close(); err != nil {
		return err
	}
	return p.MainServer.Close()
}

// btof converts a bool to a float64 (0 or 1)
func btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func main() {
	park := New()
	log.Fatal(park.Start())
}
