package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"kubepark/pkg/httptypes"

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
	mainMux.HandleFunc("/is-attraction", handleIsAttraction(config, state))
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

// Register registers the attraction with kubepark
func (a *Attraction) Register() error {
	url := a.Config.SelfURL
	if url == "" {
		url = fmt.Sprintf("http://%s:80", a.Config.Name)
	}

	registrationData := httptypes.RegisterAttractionRequest{
		URL:        url,
		BuildCost:  a.Config.BuildCost,
		RepairCost: a.Config.RepairCost,
		Size:       a.Config.Size,
	}
	data, err := json.Marshal(registrationData)
	if err != nil {
		return fmt.Errorf("failed to marshal registration data: %v", err)
	}

	resp, err := http.Post(a.Config.ParkURL+"/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to register with park: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}

	log.Printf("Successfully registered %s with park at %s", a.Config.Name, url)
	return nil
}

// Start starts both the metrics and main HTTP servers
func (a *Attraction) Start() error {
	// Register with park
	log.Printf("Registering attraction with park")
	if err := a.Register(); err != nil {
		return fmt.Errorf("failed to register attraction: %v", err)
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
				breakReq := httptypes.BreakAttractionRequest{
					URL: a.Config.SelfURL,
				}
				breakData, err := json.Marshal(breakReq)
				if err != nil {
					log.Printf("Failed to marshal break request: %v", err)
					continue
				}

				resp, err := http.Post(a.Config.ParkURL+"/break", "application/json", bytes.NewBuffer(breakData))
				if err != nil {
					log.Printf("Failed to send break request: %v", err)
					continue
				}
				resp.Body.Close()

				if resp.StatusCode == http.StatusOK {
					log.Printf("Attraction %s has broken down", a.Config.Name)
				}
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
func PayPark(w http.ResponseWriter, r *http.Request, config *Config) error {
	paymentReq := httptypes.PayParkRequest{Amount: config.Fee}
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

	Metrics.Revenue.Add(config.Fee)

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

		if config.Closed {
			Metrics.AttractionAttempts.WithLabelValues("false", "attraction_closed").Inc()
			http.Error(w, fmt.Sprintf("%s is closed", config.Name), http.StatusServiceUnavailable)
			return
		}

		// Process payment with kubepark
		if err := PayPark(w, r, config); err != nil {
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
