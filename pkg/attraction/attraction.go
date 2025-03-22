package attraction

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"kubepark/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config represents the common configuration for all attractions
type Config struct {
	Closed     bool
	Fee        float64
	VolumePath string
	ParkURL    string
	Name       string
	Duration   time.Duration
}

// Attraction represents a base attraction that can be embedded by specific attractions
type Attraction struct {
	Config        *Config
	MetricsServer *http.Server
	MainServer    *http.Server
}

// PaymentRequest represents a request to process a payment
type PaymentRequest struct {
	Amount float64 `json:"amount"`
}

// New creates a new base attraction
func New(name string, defaultFee float64, duration time.Duration) *Attraction {
	config := &Config{
		Name:     name,
		Fee:      defaultFee,
		Duration: duration,
	}

	flag.BoolVar(&config.Closed, "closed", false, "Whether the attraction is closed")
	flag.StringVar(&config.ParkURL, "park-url", "http://kubepark:80", "URL of the kubepark service")
	flag.StringVar(&config.VolumePath, "volume", "", "Path to volume for persistent storage")
	flag.Float64Var(&config.Fee, "fee", defaultFee, "Fee for using the attraction")
	flag.Parse()

	metrics.RegisterAttractionMetrics()
	metrics.Fee.Set(config.Fee)
	metrics.IsAttractionClosed.Set(btof(config.Closed))

	// Create metrics server on port 9000
	metricsServer := &http.Server{
		Addr:    ":9000",
		Handler: http.HandlerFunc(handleMetrics),
	}

	// Create main server on port 80
	mainServer := &http.Server{
		Addr:    ":80",
		Handler: http.NewServeMux(),
	}

	return &Attraction{
		Config:        config,
		MetricsServer: metricsServer,
		MainServer:    mainServer,
	}
}

// Start starts both the metrics and main HTTP servers
func (a *Attraction) Start() error {
	// Start metrics server
	go func() {
		if err := a.MetricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
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

// IsClosed returns whether the attraction is closed
func (a *Attraction) IsClosed() bool {
	return a.Config.Closed
}

// GetFee returns the attraction's fee
func (a *Attraction) GetFee() float64 {
	return a.Config.Fee
}

// GetParkURL returns the park's URL
func (a *Attraction) GetParkURL() string {
	return a.Config.ParkURL
}

// GetDuration returns the attraction's usage duration
func (a *Attraction) GetDuration() time.Duration {
	return a.Config.Duration
}

// ProcessPayment processes a payment with the park
func (a *Attraction) ProcessPayment() error {
	paymentReq := PaymentRequest{Amount: a.Config.Fee}
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

	return nil
}

// handleMetrics handles the metrics endpoint
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

// btof converts a bool to a float64 (0 or 1)
func btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
