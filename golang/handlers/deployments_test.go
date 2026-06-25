package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/client-go/kubernetes/fake"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
)

func TestDeploymentHealthHandler(t *testing.T) {

	replicas := int32(1)

	clientset := fake.NewSimpleClientset(
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
			},
			Status: appsv1.DeploymentStatus{
				ReadyReplicas: 1,
			},
		},
	)

	handler := DeploymentHealthHandler(clientset)

	req := httptest.NewRequest(http.MethodGet, "/deployments/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.Contains(t, string(body), "test")
}
func TestReadyzHandler(t *testing.T) {

	clientset := fake.NewSimpleClientset()

	handler := K8sHealthHandler(clientset)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}