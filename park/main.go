package main

import (
	"context"
	"fmt"
	"io"
	"kubepark/pkg/k8s"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Park represents the amusement park simulator
type Park struct {
	Config        *Config
	MetricsServer *http.Server
	MainServer    *http.Server
	State         *StateManager
	GuestManager  *GuestJobManager
}

// New creates a new park simulator
func New() *Park {
	config := &Config{}

	RegisterFlags(config)

	// Check for existing park service
	if err := checkForExistingPark(); err != nil {
		log.Fatalf("Failed park check: %v", err)
	}

	// Initialize game state
	state, err := NewStateManager(config)
	if err != nil {
		log.Fatalf("Failed to initialize game state: %v", err)
	}

	// Initialize guest job manager
	guestManager, err := NewGuestJobManager()
	if err != nil {
		log.Fatalf("Failed to initialize guest job manager: %v", err)
	}

	metrics.EntranceFee.Set(config.EntranceFee)
	metrics.IsParkClosed.Set(btof(config.Closed))
	metrics.OpensAt.Set(float64(config.OpensAt))
	metrics.ClosesAt.Set(float64(config.ClosesAt))

	// Create metrics server on port 9000
	r := prometheus.NewRegistry()
	RegisterParkMetrics(r)
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	metricsServer := &http.Server{
		Addr:    ":9000",
		Handler: metricsMux,
	}

	// Create main server on port 80
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/park-status", handleParkStatus(state))
	mainMux.HandleFunc("/pay", handlePayPark(state))
	mainMux.HandleFunc("/enter", handleEnterPark(state))
	mainServer := &http.Server{
		Addr:    ":80",
		Handler: mainMux,
	}

	return &Park{
		Config:        config,
		MetricsServer: metricsServer,
		MainServer:    mainServer,
		State:         state,
		GuestManager:  guestManager,
	}
}

// Start starts the park simulator
func (p *Park) Start() error {
	// Start metrics server
	go func() {
		log.Printf("Starting metrics server on port 9000")
		if err := p.MetricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Start the park simulation loop
	go func() {
		log.Printf("Starting park simulation loop")
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		ctx := context.Background()

		for range ticker.C {
			// Speed up simulation time
			p.State.SetTime(p.State.GetTime().Add(time.Second * 100))

			// Update metrics
			metrics.Time.Set(float64(p.State.GetTime().Unix()))
			metrics.Money.Set(p.State.GetMoney())

			currentHour := p.State.GetTime().Hour()
			isClosed := currentHour < p.Config.OpensAt || currentHour >= p.Config.ClosesAt

			if isClosed {
				if err := p.GuestManager.CleanupJobs(ctx); err != nil {
					log.Printf("Failed to cleanup all jobs during closed hours: %v", err)
				}

				continue
			}

			// Randomly decide whether to create a new guest (10% chance)
			if rand.Float64() < 0.1 {
				// Create new guests if park is open and not at capacity
				url := p.Config.SelfURL
				if url == "" {
					url = "http://park:80"
				}

				if err := p.GuestManager.CreateGuestJob(ctx, p.Config.Image, url); err != nil {
					log.Printf("Failed to create guest job: %v", err)
				}
			}
		}
	}()

	// Start main server
	log.Printf("Starting main server on port 80")
	return p.MainServer.ListenAndServe()
}

// Stop gracefully stops the park simulator
func (p *Park) Stop() error {
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

// checkForExistingPark checks if another park service exists in the cluster
func checkForExistingPark() error {
	var parks []interface{}

	decoder := func(r io.Reader, v *[]interface{}, ip string) error {
		*v = append(*v, true)
		return nil
	}

	err := k8s.DiscoverServices("/park-status", &parks, decoder)
	if err != nil {
		return err
	}

	for len(parks) > 0 {
		return fmt.Errorf("another park service is already running in the cluster")
	}

	return nil
}
