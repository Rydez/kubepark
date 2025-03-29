package k8s

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func DiscoverServices(endpoint string, v *[]interface{}, decoder func(r io.Reader, v *[]interface{}, ip string) error) error {
	// Get in-cluster config
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to get in-cluster config: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %v", err)
	}

	// Get all pods in the cluster
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}

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
