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

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Attraction represents a base attraction that can be embedded by specific attractions
type Attraction struct {
	Config        *Config
	MetricsServer *http.Server
	MainServer    *http.Server
}

// New creates a new base attraction
func New(config *Config, defaultFee float64) *Attraction {
	RegisterFlags(config, defaultFee)

	RegisterAttractionMetrics()
	Metrics.Fee.Set(config.Fee)
	Metrics.IsAttractionClosed.Set(btof(config.Closed))

	// Create metrics server on port 9000
	metricsMux := http.NewServeMux()
	metricsMux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	metricsServer := &http.Server{
		Addr:    ":9000",
		Handler: metricsMux,
	}

	// Create main server on port 80
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/status", handleStatus())
	mainServer := &http.Server{
		Addr:    ":80",
		Handler: mainMux,
	}

	return &Attraction{
		Config:        config,
		MetricsServer: metricsServer,
		MainServer:    mainServer,
	}
}

// handleStatus handles the status endpoint
func handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
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
		return fmt.Errorf("failed to register with kubepark: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}

	log.Printf("Successfully registered %s with kubepark at %s", a.Config.Name, url)
	return nil
}

// Start starts both the metrics and main HTTP servers
func (a *Attraction) Start() error {
	// Register with kubepark
	if err := a.Register(); err != nil {
		return fmt.Errorf("failed to register attraction: %v", err)
	}

	// Start metrics server
	go func() {
		if err := a.MetricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Start the attraction simulation loop
	go func() {
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
func (a *Attraction) PayPark() error {
	paymentReq := httptypes.PayParkRequest{Amount: a.Config.Fee}
	paymentData, err := json.Marshal(paymentReq)
	if err != nil {
		return err
	}

	resp, err := http.Post(a.Config.ParkURL+"/pay", "application/json", bytes.NewBuffer(paymentData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment failed with status: %d", resp.StatusCode)
	}

	Metrics.Revenue.Add(a.Config.Fee)

	return nil
}

// ValidatePayment validates if a guest has enough money to use the attraction
func (a *Attraction) ValidatePayment(w http.ResponseWriter, r *http.Request) (float64, error) {
	var useReq httptypes.UseAttractionRequest
	if err := json.NewDecoder(r.Body).Decode(&useReq); err != nil {
		return 0, fmt.Errorf("invalid payment request: %v", err)
	}

	if useReq.GuestMoney < a.Config.Fee {
		return 0, fmt.Errorf("insufficient funds. Fee is $%.2f but guest has $%.2f", a.Config.Fee, useReq.GuestMoney)
	}

	return useReq.GuestMoney, nil
}

// btof converts a bool to a float64 (0 or 1)
func btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
