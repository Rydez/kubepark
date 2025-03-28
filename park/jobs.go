package main

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GuestJobManager manages Kubernetes jobs for park guests
type GuestJobManager struct {
	clientset *kubernetes.Clientset
}

// NewGuestJobManager creates a new guest job manager
func NewGuestJobManager() (*GuestJobManager, error) {
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
	}, nil
}

// CreateGuestJob creates a new guest job
func (m *GuestJobManager) CreateGuestJob(ctx context.Context, image string, parkURL string) error {
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
							Image:   image,
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

	_, err := m.clientset.BatchV1().Jobs("guests").Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create guest job: %w", err)
	}

	return nil
}

// CleanupJobs removes all jobs
func (m *GuestJobManager) CleanupJobs(ctx context.Context) error {
	jobs, err := m.clientset.BatchV1().Jobs("guests").List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	backgroundDeletion := metav1.DeletePropagationBackground
	for _, job := range jobs.Items {
		err := m.clientset.BatchV1().Jobs("guests").Delete(ctx, job.Name, metav1.DeleteOptions{
			PropagationPolicy: &backgroundDeletion,
		})
		if err != nil {
			return fmt.Errorf("failed to delete job %s: %w", job.Name, err)
		}
	}

	return nil
}

// int32Ptr returns a pointer to the given int32 value
func int32Ptr(i int32) *int32 {
	return &i
}
