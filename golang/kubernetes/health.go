package kubernetes

import (
	"time"

	"k8s.io/client-go/kubernetes"
)

type ClusterHealth struct {
	Connected     bool   `json:"connected"`
	ServerVersion string `json:"serverVersion,omitempty"`
	LatencyMs     int64  `json:"latencyMs,omitempty"`
	Error         string `json:"error,omitempty"`
}

func CheckAPIServer(clientset kubernetes.Interface) ClusterHealth {

	start := time.Now()

	version, err := clientset.Discovery().ServerVersion()
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return ClusterHealth{
			Connected: false,
			Error:     err.Error(),
		}
	}

	return ClusterHealth{
		Connected:     true,
		ServerVersion: version.String(),
		LatencyMs:     latency,
	}
}