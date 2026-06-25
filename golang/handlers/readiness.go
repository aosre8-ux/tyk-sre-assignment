package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TykTechnologies/tyk-sre-assignment/kubernetes"
	k8sclient "k8s.io/client-go/kubernetes"
)

func K8sHealthHandler(clientset k8sclient.Interface) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		result := kubernetes.CheckAPIServer(clientset)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}