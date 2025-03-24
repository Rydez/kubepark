package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Park represents the amusement park simulator
type Park struct {
	Config        *Config
	MetricsServer *http.Server
	MainServer    *http.Server
	State         *GameState
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
	gameState, err := NewGameState(config)
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
	mainMux.HandleFunc("/is-park", handleIsPark())
	mainMux.HandleFunc("/pay", handlePayPark(gameState))
	mainMux.HandleFunc("/enter", handleEnterPark(gameState))
	mainMux.HandleFunc("/attractions", handleListAttractions(gameState))
	mainMux.HandleFunc("/register", handleRegisterAttraction(gameState))
	mainMux.HandleFunc("/break", handleBreakAttraction(gameState))
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
			p.State.SetTime(p.State.GetTime().Add(time.Second * 10))

			// Update metrics
			metrics.Time.Set(float64(p.State.GetTime().Unix()))
			metrics.Money.Set(p.State.GetMoney())

			// Save state every minute
			if time.Since(p.State.GetTime()) > time.Minute {
				if err := p.State.Save(); err != nil {
					log.Printf("Failed to save state: %v", err)
				}
			}

			// Check attraction statuses
			p.checkAttractionPending()

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

				if err := p.GuestManager.CreateGuestJob(ctx, url); err != nil {
					log.Printf("Failed to create guest job: %v", err)
				}
			}
		}
	}()

	// Start the main server
	log.Printf("Starting main server on port 80")
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

// checkForExistingPark checks if another park service exists in the cluster
func checkForExistingPark() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	// Look for pods with port 80 across all namespaces
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	// Count running park services
	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}

		// Skip the current pod
		if pod.Name == os.Getenv("HOSTNAME") {
			continue
		}

		// Check if pod has port 80
		hasPort80 := false
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				if port.ContainerPort == 80 {
					hasPort80 = true
					break
				}
			}
			if hasPort80 {
				break
			}
		}

		if !hasPort80 {
			continue
		}

		// Try to connect to the /is-park endpoint
		client := &http.Client{
			Timeout: time.Second * 2,
		}

		// Try to connect to the pod's IP
		podIP := pod.Status.PodIP
		if podIP == "" {
			continue
		}

		resp, err := client.Get(fmt.Sprintf("http://%s/is-park", podIP))
		if err != nil {
			continue // Skip if we can't connect
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return fmt.Errorf("another park service is already running in the cluster")
		}
	}

	return nil
}

// checkAttractionPending checks the status of all attractions and updates their pending state
func (p *Park) checkAttractionPending() {
	client := &http.Client{
		Timeout: time.Second * 2,
	}

	for _, attraction := range p.State.GetAttractions() {
		// Skip if we can't reach the attraction
		resp, err := client.Get(attraction.URL + "/status")
		if err != nil {
			// If we can't reach the attraction and it's not pending, mark it as pending
			if !attraction.IsPending {
				if err := p.State.SetAttractionPending(attraction.URL, true); err != nil {
					log.Printf("Failed to update attraction status: %v", err)
				}
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// If we get an error response and it's not pending, mark it as pending
			if !attraction.IsPending {
				if err := p.State.SetAttractionPending(attraction.URL, true); err != nil {
					log.Printf("Failed to update attraction status: %v", err)
				}
			}
			continue
		}

		// Update the attraction's pending status
		if err := p.State.SetAttractionPending(attraction.URL, false); err != nil {
			log.Printf("Failed to update attraction status: %v", err)
		}
	}
}
