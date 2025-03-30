package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kubepark/pkg/httptypes"
	"net/http"
	"os"
	"time"

	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewClient() (*kubernetes.Clientset, error) {
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

	return clientset, nil
}

func DiscoverServices(endpoint string, v *[]interface{}, decoder func(r io.Reader, v *[]interface{}, ip string) error) error {
	clientset, err := NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	// Get all pods in the cluster
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	// Check each pod for the specified endpoint
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
		resp, err := client.Get(fmt.Sprintf("http://%s/%s", podIP, endpoint))
		if err != nil {
			continue // Skip if we can't connect
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			decoder(resp.Body, v, podIP)
		}
	}

	return nil
}

func DiscoverGuests() (*v1.JobList, error) {
	clientset, err := NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	jobs, err := clientset.BatchV1().Jobs("guests").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "status.active=1",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %v", err)
	}

	return jobs, nil
}

func DiscoverAttractions() ([]httptypes.Attraction, error) {
	var attractions []interface{}

	decoder := func(r io.Reader, v *[]interface{}, ip string) error {
		var attraction httptypes.Attraction
		if err := json.NewDecoder(r).Decode(&attraction); err != nil {
			return err
		}

		*v = append(*v, httptypes.Attraction{
			URL:  fmt.Sprintf("http://%s", ip),
			Fee:  attraction.Fee,
			Size: attraction.Size,
		})

		return nil
	}

	err := DiscoverServices("/attraction-status", &attractions, decoder)

	if err != nil {
		return nil, err
	}

	var typedAttractions []httptypes.Attraction
	for _, a := range attractions {
		if attr, ok := a.(httptypes.Attraction); ok {
			typedAttractions = append(typedAttractions, attr)
		}
	}

	return typedAttractions, nil
}
