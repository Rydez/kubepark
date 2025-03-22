package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"kubepark/pkg/state"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Park represents the amusement park simulator
type Park struct {
	Config        *Config
	MetricsServer *http.Server
	MainServer    *http.Server
	State         *state.GameState
	GuestManager  *GuestJobManager
}

// New creates a new park simulator
func New() *Park {
	config := &Config{}

	RegisterFlags(config)

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

	// Initialize guest job manager
	guestManager, err := NewGuestJobManager(config.Namespace)
	if err != nil {
		log.Fatalf("Failed to initialize guest job manager: %v", err)
	}

	RegisterParkMetrics()
	metrics.EntranceFee.Set(config.EntranceFee)
	metrics.IsParkClosed.Set(btof(config.Closed))
	metrics.OpensAt.Set(float64(config.OpensAt))
	metrics.ClosesAt.Set(float64(config.ClosesAt))

	// Create metrics server on port 9000
	metricsMux := http.NewServeMux()
	metricsMux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	metricsServer := &http.Server{
		Addr:    ":9000",
		Handler: metricsMux,
	}

	// Create main server on port 80
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/pay", handlePayPark(gameState))
	mainMux.HandleFunc("/enter", handleEnterPark(gameState))
	mainMux.HandleFunc("/attractions", handleListAttractions(gameState))
	mainMux.HandleFunc("/register", handleRegisterAttraction(gameState))
	mainServer := &http.Server{
		Addr:    ":80",
		Handler: mainMux,
	}

	return &Park{
		Config:        config,
		MetricsServer: metricsServer,
		MainServer:    mainServer,
		State:         gameState,
		GuestManager:  guestManager,
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

		ctx := context.Background()

		for range ticker.C {
			// Update metrics
			metrics.Money.Set(p.State.GetMoney())
			metrics.Time.Set(float64(p.State.GetTime().Unix()))

			// Save state every minute
			if time.Since(p.State.GetTime()) > time.Minute {
				if err := p.State.Save(); err != nil {
					log.Printf("Failed to save state: %v", err)
				}
			}

			// Create new guests if park is open and not at capacity
			if !p.State.IsClosed() {
				// TODO: Add capacity check based on mode
				if err := p.GuestManager.CreateGuestJob(ctx, "http://kubepark:80"); err != nil {
					log.Printf("Failed to create guest job: %v", err)
				}
			}

			// Cleanup old guest jobs
			if err := p.GuestManager.CleanupOldJobs(ctx); err != nil {
				log.Printf("Failed to cleanup old jobs: %v", err)
			}
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
