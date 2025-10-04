package main

import (
	"context"
	"fmt"
	"io"
	"kubepark/pkg/k8s"
	"kubepark/pkg/logger"
	"log/slog"
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
	GrafanaLive   *GrafanaLiveClient
}

// New creates a new park simulator
func New() *Park {
	config := &Config{}

	RegisterFlags(config)

	// Initialize logger with configured level
	logger.InitLogger(config.LogLevel)

	// Check for existing park service
	if err := checkForExistingPark(); err != nil {
		slog.Error("Failed park check", "error", err)
		panic(err)
	}

	// Initialize game state
	state, err := NewStateManager(config)
	if err != nil {
		slog.Error("Failed to initialize game state", "error", err)
		panic(err)
	}

	// Initialize guest job manager
	guestManager, err := NewGuestJobManager()
	if err != nil {
		slog.Error("Failed to initialize guest job manager", "error", err)
		panic(err)
	}

	// Initialize Grafana Live client
	grafanaLive := NewGrafanaLiveClient(config.GrafanaURL, config.GrafanaAPIKey)

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
	mainMux.HandleFunc("/park-status", handleStatus(config, state))
	mainMux.HandleFunc("/transaction", handleTransaction(state))
	mainMux.HandleFunc("/enter", handleEnter(state))
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
		GrafanaLive:   grafanaLive,
	}
}

// Start starts the park simulator
func (p *Park) Start() error {
	// Start metrics server
	go func() {
		slog.Info("Starting metrics server on port 9000")
		if err := p.MetricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Metrics server failed", "error", err)
			panic(err)
		}
	}()

	// Start the park simulation loop
	go func() {
		slog.Info("Starting park simulation loop")
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		ctx := context.Background()

		for range ticker.C {
			time := p.State.GetTime().Add(time.Second * 100)

			// Speed up simulation time
			err := p.State.SetTime(time)
			if err != nil {
				slog.Error("Failed to set time", "error", err)
			}

			// Update metrics
			metrics.Time.Set(float64(time.Unix()))
			metrics.Money.Set(p.State.GetMoney())

			// Push to Grafana Live
			if err := p.GrafanaLive.PushMetric("park_time", float64(time.Unix()*1000), nil); err != nil {
				slog.Warn("Failed to push time metric to Grafana Live", "error", err)
			}
			if err := p.GrafanaLive.PushMetric("park_money", p.State.GetMoney(), nil); err != nil {
				slog.Warn("Failed to push money metric to Grafana Live", "error", err)
			}

			if isClosed(p.Config, time) {
				foundJobs, err := p.GuestManager.CleanupJobs(ctx)
				if err != nil {
					slog.Error("Failed to cleanup all jobs during closed hours", "error", err)
				}

				if foundJobs > 0 {
					slog.Info("Removed guests after park closed", "count", foundJobs)
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
					slog.Warn("Failed to create guest job", "error", err)
				}
			}
		}
	}()

	// Start main server
	slog.Info("Starting main server on port 80")
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
	if err := park.Start(); err != nil {
		slog.Error("Park failed to start", "error", err)
		panic(err)
	}
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

func isClosed(config *Config, time time.Time) bool {
	hour := time.Hour()
	return config.Closed || hour < config.OpensAt || hour >= config.ClosesAt
}
