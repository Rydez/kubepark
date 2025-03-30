package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"kubepark/pkg/constants"
	"kubepark/pkg/httptypes"
	"kubepark/pkg/k8s"

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

	// Initialize state manager
	state, err := NewStateManager(config.VolumePath)
	if err != nil {
		log.Fatalf("Failed to initialize state manager: %v", err)
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

		if err := PayPark(a.Config, a.Config.RepairCost); err != nil {
			return fmt.Errorf("failed to pay for repair: %v", err)
		}
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

	jobs, err := k8s.DiscoverGuests()
	if err != nil {
		return fmt.Errorf("failed to discover guest jobs: %v", err)
	}

	usedSpace += constants.GuestSize * float64(len(jobs.Items))

	if park.TotalSpace-usedSpace < a.Config.Size {
		return fmt.Errorf("not enough space in the park")
	}

	if err := PayPark(a.Config, a.Config.BuildCost); err != nil {
		return fmt.Errorf("failed to pay for build: %v", err)
	}

	log.Printf("Successfully started %s", a.Config.Name)
	return nil
}

// Start starts both the metrics and main HTTP servers
func (a *Attraction) Start() error {
	// Register with park
	log.Printf("Checking if attraction can start")
	if err := a.BeforeStart(); err != nil {
		return fmt.Errorf("failed check to see if attraction can start: %v", err)
	}

	// Start metrics server
	go func() {
		log.Printf("Starting metrics server on port 9000")
		if err := a.MetricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Start the attraction simulation loop
	go func() {
		log.Printf("Starting attraction simulation loop")
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Random chance to break the attraction (1% chance per second)
			if rand.Float64() < 0.01 {
				a.State.SetBroken(true)
			}
		}
	}()

	// Start main server
	log.Printf("Starting main server on port 80")
	return a.MainServer.ListenAndServe()
}

// Stop gracefully stops both HTTP servers
func (a *Attraction) Stop() error {
	if err := a.MetricsServer.Close(); err != nil {
		return err
	}
	return a.MainServer.Close()
}

// PayPark processes a payment with the park
func PayPark(config *Config, amount float64) error {
	paymentReq := httptypes.PayParkRequest{Amount: amount}
	paymentData, err := json.Marshal(paymentReq)
	if err != nil {
		return err
	}

	resp, err := http.Post(config.ParkURL+"/pay", "application/json", bytes.NewBuffer(paymentData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment failed with status: %d", resp.StatusCode)
	}

	Metrics.Revenue.Add(amount)

	return nil
}

// btof converts a bool to a float64 (0 or 1)
func btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
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
		if err := PayPark(config, config.Fee); err != nil {
			log.Printf("Failed to process payment: %v", err)
			Metrics.AttractionAttempts.WithLabelValues("false", "payment_failed").Inc()
			http.Error(w, "Payment failed", http.StatusInternalServerError)
			return
		}

		// Simulate usage duration
		time.Sleep(config.Duration)

		// Call after use hook if set
		if afterUse != nil {
			if err := afterUse(); err != nil {
				log.Printf("After use hook failed: %v", err)
				Metrics.AttractionAttempts.WithLabelValues("false", "hook_failed").Inc()
				http.Error(w, "Failed to cleanup attraction", http.StatusInternalServerError)
				return
			}
		}

		Metrics.AttractionAttempts.WithLabelValues("true", "success").Inc()
		w.WriteHeader(http.StatusOK)
	}
}
