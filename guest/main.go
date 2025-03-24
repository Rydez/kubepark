package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"kubepark/pkg/httptypes"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Guest metrics
	MoneySpent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "money_spent",
		Help: "Total money spent by the guest",
	})

	AttractionsVisited = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "attractions_visited",
		Help: "Number of attractions visited",
	})

	// Configuration
	config struct {
		ParkURL   string
		Money     float64
		StartTime time.Time
		EndTime   time.Time
	}
)

// ParkResponse represents a response from the park
type ParkResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Attraction represents an attraction in the park
type Attraction struct {
	Name string  `json:"name"`
	Fee  float64 `json:"fee"`
}

func main() {
	// Register metrics
	r := prometheus.NewRegistry()
	r.MustRegister(MoneySpent)
	r.MustRegister(AttractionsVisited)

	// Parse flags
	flag.StringVar(&config.ParkURL, "park-url", "http://kubepark:80", "URL of the kubepark service")
	flag.Parse()

	// Set fixed values
	config.Money = 100 // Each guest starts with $100

	// Set start and end times for the day
	now := time.Now()
	config.StartTime = time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
	config.EndTime = time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, now.Location())

	// Start metrics server
	go func() {
		log.Printf("Starting metrics server on port 9000")
		http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
		if err := http.ListenAndServe(":9000", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Try to enter the park
	log.Printf("Entering park")
	if err := enterPark(); err != nil {
		log.Printf("Failed to enter park: %v", err)
		return
	}

	// Start exploring attractions
	log.Printf("Starting attraction loop")
	for {
		// Check if park is still open
		if time.Now().After(config.EndTime) {
			log.Println("Park is closed, leaving")
			break
		}

		// Visit a random attraction
		if err := visitAttraction(); err != nil {
			log.Printf("Failed to visit attraction: %v", err)
		}

		// Random chance (30%) that guest decides to leave early
		if rand.Float64() < 0.30 {
			log.Println("Guest decided to leave early")
			break
		}

		// Take a break between attractions
		time.Sleep(time.Duration(rand.Intn(30)+30) * time.Second)
	}

	log.Printf("Guest finished their visit. Attractions visited: %d", AttractionsVisited)
}

func enterPark() error {
	// Make request to enter park
	resp, err := http.Post(config.ParkURL+"/enter", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var parkResp ParkResponse
	if err := json.NewDecoder(resp.Body).Decode(&parkResp); err != nil {
		return err
	}

	if !parkResp.Success {
		return fmt.Errorf("failed to enter park: %s", parkResp.Message)
	}

	log.Println("Successfully entered the park")
	return nil
}

func visitAttraction() error {
	// Get list of available attractions
	resp, err := http.Get(config.ParkURL + "/attractions")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var attractions []Attraction
	if err := json.NewDecoder(resp.Body).Decode(&attractions); err != nil {
		return err
	}

	if len(attractions) == 0 {
		return fmt.Errorf("no attractions available")
	}

	// Choose a random attraction
	randAttraction := attractions[rand.Intn(len(attractions))]

	// Create use request with remaining money
	useReq := httptypes.UseAttractionRequest{
		GuestMoney: config.Money,
	}
	useData, err := json.Marshal(useReq)
	if err != nil {
		return fmt.Errorf("failed to marshal use request: %v", err)
	}

	// Visit the attraction
	resp, err = http.Post(fmt.Sprintf("%s/use", randAttraction.Name), "application/json", bytes.NewBuffer(useData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to use attraction %s: %s", randAttraction.Name, resp.Status)
	}

	// Update metrics and money
	MoneySpent.Add(randAttraction.Fee)
	AttractionsVisited.Inc()
	config.Money -= randAttraction.Fee

	log.Printf("Visited %s for $%.2f", randAttraction.Name, randAttraction.Fee)
	return nil
}
