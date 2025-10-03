package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kubepark/pkg/httptypes"
	"kubepark/pkg/k8s"
	"kubepark/pkg/logger"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Attraction represents a base attraction that can be embedded by specific attractions
type Attraction struct {
	Config        *Config
	MetricsServer *http.Server
	MainServer    *http.Server
	State         *StateManager
}

// New creates a new base attraction
func New(config *Config, defaultFee float64, afterUse func() error) *Attraction {
	RegisterFlags(config, defaultFee)

	// Initialize logger with configured level
	logger.InitLogger(config.LogLevel)

	// Initialize state manager
	state, err := NewStateManager(config.VolumePath)
	if err != nil {
		slog.Error("Failed to initialize state manager", "error", err)
		panic(err)
	}

	r := prometheus.NewRegistry()
	RegisterAttractionMetrics(r)
	Metrics.Fee.Set(config.Fee)
	Metrics.IsAttractionClosed.Set(btof(config.Closed))

	// Create metrics server on port 9000
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	metricsServer := &http.Server{
		Addr:    ":9000",
		Handler: metricsMux,
	}

	// Create main server on port 80
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/use", handleUse(config, state, afterUse))
	mainMux.HandleFunc("/attraction-status", handleAttractionStatus(config, state))
	mainServer := &http.Server{
		Addr:    ":80",
		Handler: mainMux,
	}

	return &Attraction{
		Config:        config,
		MetricsServer: metricsServer,
		MainServer:    mainServer,
		State:         state,
	}
}

// BeforeStart checks if there's enough space in the park and enough money to build the attraction.
func (a *Attraction) BeforeStart() error {
	resp, err := http.Get(a.Config.ParkURL + "/park-status")
	if err != nil {
		return fmt.Errorf("failed to get park status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("park status check failed with status: %d", resp.StatusCode)
	}

	var park httptypes.Park
	if err := json.NewDecoder(resp.Body).Decode(&park); err != nil {
		return fmt.Errorf("failed to decode park status: %v", err)
	}

	if a.State.IsBroken() {
		if a.Config.RepairCost > park.Money {
			return fmt.Errorf("not enough money to repair attraction")
		}

		if err := ParkTransaction(a.Config, -a.Config.RepairCost); err != nil {
			return fmt.Errorf("failed to pay for repair: %v", err)
		}

		err = a.State.SetBroken(false)
		if err != nil {
			return fmt.Errorf("failed to set attraction not broken: %v", err)
		}

		slog.Info("Successfully repaired attraction", "name", a.Config.Name)
		return nil
	}

	if !park.IsClosed {
		return fmt.Errorf("cannot build attraction when park is open")
	}

	if a.Config.BuildCost > park.Money {
		return fmt.Errorf("not enough money to build attraction")
	}

	attractions, err := k8s.DiscoverAttractions()
	if err != nil {
		return fmt.Errorf("failed to discover attractions: %v", err)
	}

	usedSpace := 0.0
	for _, attraction := range attractions {
		usedSpace += attraction.Size
	}

	if park.TotalSpace-usedSpace < a.Config.Size {
		return fmt.Errorf("not enough space in the park")
	}

	if err := ParkTransaction(a.Config, -a.Config.BuildCost); err != nil {
		return fmt.Errorf("failed to pay for build: %v", err)
	}

	slog.Info("Successfully built new attraction", "name", a.Config.Name)
	return nil
}

// Start starts both the metrics and main HTTP servers
func (a *Attraction) Start() error {
	// Register with park
	slog.Info("Checking if attraction can start")
	if err := a.BeforeStart(); err != nil {
		return fmt.Errorf("failed check to see if attraction can start: %v", err)
	}

	// Start metrics server
	go func() {
		slog.Info("Starting metrics server on port 9000")
		if err := a.MetricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Metrics server failed", "error", err)
			panic(err)
		}
	}()

	// Start the attraction simulation loop
	go func() {
		slog.Info("Starting attraction simulation loop")
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Random chance to break the attraction (0.1% chance per second)
			if !a.State.IsBroken() && rand.Float64() < 0.001 {
				slog.Info("Attraction has broken down", "name", a.Config.Name)
				err := a.State.SetBroken(true)
				if err != nil {
					slog.Error("Failed to set attraction broken", "error", err)
				}
			}
		}
	}()

	// Start main server
	slog.Info("Starting main server on port 80")
	return a.MainServer.ListenAndServe()
}

// Stop gracefully stops both HTTP servers
func (a *Attraction) Stop() error {
	if err := a.MetricsServer.Close(); err != nil {
		return err
	}
	return a.MainServer.Close()
}

// ParkTransaction processes a transaction with the park
func ParkTransaction(config *Config, amount float64) error {
	req := httptypes.TransactionRequest{Amount: amount}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(config.ParkURL+"/transaction", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment failed with status: %d", resp.StatusCode)
	}

	if amount > 0 {
		Metrics.Revenue.Add(amount)
	} else {
		Metrics.Costs.Add(-amount)
	}

	return nil
}

// btof converts a bool to a float64 (0 or 1)
func btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
