package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"kubepark/pkg/httptypes"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
		ParkURL string
		Money   float64
	}
)

// Attraction represents an attraction in the park
type Attraction struct {
	URL string  `json:"url"`
	Fee float64 `json:"fee"`
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

	log.Printf("Guest finished their visit.")
}

func enterPark() error {
	// Make request to enter park
	resp, err := http.Post(config.ParkURL+"/enter", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to enter park: %s", resp.Status)
	}

	log.Println("Successfully entered the park")
	return nil
}

func visitAttraction() error {
	// Get list of available attractions using Kubernetes API
	attractions, err := discoverAttractions()
	if err != nil {
		return fmt.Errorf("failed to discover attractions: %v", err)
	}

	if len(attractions) == 0 {
		return fmt.Errorf("no attractions available")
	}

	// Choose a random attraction
	randAttraction := attractions[rand.Intn(len(attractions))]

	// Check if guest has enough money
	if config.Money < randAttraction.Fee {
		return fmt.Errorf("insufficient funds. Fee is $%.2f but guest has $%.2f", randAttraction.Fee, config.Money)
	}

	// Visit the attraction
	resp, err := http.Post(fmt.Sprintf("%s/use", randAttraction.URL), "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to use attraction %s: %s", randAttraction.URL, resp.Status)
	}

	// Update metrics and money
	MoneySpent.Add(randAttraction.Fee)
	AttractionsVisited.Inc()
	config.Money -= randAttraction.Fee

	log.Printf("Visited %s for $%.2f", randAttraction.URL, randAttraction.Fee)
	return nil
}

func discoverAttractions() ([]Attraction, error) {
	// Get in-cluster config
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	// Get all pods in the cluster
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	var attractions []Attraction
	client := &http.Client{
		Timeout: time.Second * 2,
	}

	// Check each pod for the /is-attraction endpoint
	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}

		// Skip the current pod
		if pod.Name == os.Getenv("HOSTNAME") {
			continue
		}

		// Try to connect to the pod's IP
		podIP := pod.Status.PodIP
		if podIP == "" {
			continue
		}

		// Check if this is an attraction
		resp, err := client.Get(fmt.Sprintf("http://%s/is-attraction", podIP))
		if err != nil {
			continue // Skip if we can't connect
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var feeResponse httptypes.IsAttractionResponse
			if err := json.NewDecoder(resp.Body).Decode(&feeResponse); err != nil {
				log.Printf("Failed to decode fee response: %v", err)
				continue
			}

			attractions = append(attractions, Attraction{
				URL: fmt.Sprintf("http://%s", podIP),
				Fee: feeResponse.Fee,
			})
		}
	}

	return attractions, nil
}
