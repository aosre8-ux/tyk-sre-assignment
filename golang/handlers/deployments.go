package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TykTechnologies/tyk-sre-assignment/kubernetes"
	clientgoset "k8s.io/client-go/kubernetes"
)

func DeploymentHealthHandler(
	clientset clientgoset.Interface,
) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		deployments, err := kubernetes.GetDeploymentHealth(clientset)

		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error fetching deployments: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)
                pretty := r.URL.Query().Get("pretty") == "true"

        if pretty {
            encoder := json.NewEncoder(w)
            encoder.SetIndent("", "  ")
            _ = encoder.Encode(deployments)
            return
        }
		json.NewEncoder(w).Encode(deployments)
	}
}
