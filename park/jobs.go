package main

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GuestJobManager manages Kubernetes jobs for park guests
type GuestJobManager struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewGuestJobManager creates a new guest job manager
func NewGuestJobManager(namespace string) (*GuestJobManager, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &GuestJobManager{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

// CreateGuestJob creates a new guest job
func (m *GuestJobManager) CreateGuestJob(ctx context.Context, parkURL string) error {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "guest-",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "guest",
							Image:   "kubepark:latest",
							Command: []string{"/opt/kubepark/internal/guest"},
							Args: []string{
								"--park-url", parkURL,
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: int32Ptr(4),
		},
	}

	_, err := m.clientset.BatchV1().Jobs(m.namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create guest job: %w", err)
	}

	return nil
}

// CleanupOldJobs removes completed or failed jobs older than 1 hour
func (m *GuestJobManager) CleanupOldJobs(ctx context.Context) error {
	jobs, err := m.clientset.BatchV1().Jobs(m.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	cutoff := time.Now().Add(-1 * time.Hour)
	for _, job := range jobs.Items {
		if job.Status.Succeeded == 1 || job.Status.Failed == 1 {
			if job.Status.CompletionTime != nil && job.Status.CompletionTime.Time.Before(cutoff) {
				err := m.clientset.BatchV1().Jobs(m.namespace).Delete(ctx, job.Name, metav1.DeleteOptions{})
				if err != nil {
					return fmt.Errorf("failed to delete job %s: %w", job.Name, err)
				}
			}
		}
	}

	return nil
}

// int32Ptr returns a pointer to the given int32 value
func int32Ptr(i int32) *int32 {
	return &i
}
