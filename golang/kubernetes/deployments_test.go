package kubernetes

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
)

func TestGetDeploymentHealth(t *testing.T) {

	replicas := int32(2)

	clientset := fake.NewSimpleClientset(
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-deploy",
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
			},
			Status: appsv1.DeploymentStatus{
				ReadyReplicas: 2,
			},
		},
	)

	result, err := GetDeploymentHealth(clientset)

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	assert.Equal(t, "test-deploy", result[0].Name)
	assert.Equal(t, int32(2), result[0].Desired)
	assert.Equal(t, int32(2), result[0].Ready)
	assert.True(t, result[0].Healthy)
}

func TestGetDeploymentHealth_Unhealthy(t *testing.T) {

	replicas := int32(3)

	clientset := fake.NewSimpleClientset(
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bad-deploy",
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

	result, err := GetDeploymentHealth(clientset)

	assert.NoError(t, err)
	assert.False(t, result[0].Healthy)
}