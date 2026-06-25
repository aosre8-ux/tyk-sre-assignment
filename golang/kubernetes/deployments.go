package kubernetes

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentStatus struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Desired   int32  `json:"desired"`
	Ready     int32  `json:"ready"`
	Healthy   bool   `json:"healthy"`
}

func GetDeploymentHealth(clientset kubernetes.Interface) ([]DeploymentStatus, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deployments, err := clientset.
		AppsV1().
		Deployments("").
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	result := make([]DeploymentStatus, 0, len(deployments.Items))

	for _, d := range deployments.Items {

		desired := int32(0)
		if d.Spec.Replicas != nil {
			desired = *d.Spec.Replicas
		}

		ready := d.Status.ReadyReplicas

		result = append(result, DeploymentStatus{
			Namespace: d.Namespace,
			Name:      d.Name,
			Desired:   desired,
			Ready:     ready,
			Healthy:   desired == ready && desired > 0,
		})
	}

	return result, nil
}
