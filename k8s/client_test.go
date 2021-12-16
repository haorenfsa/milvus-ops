package k8s

import (
	"context"
	"testing"

	"github.com/haorenfsa/milvus-ops/storage/file"
	"github.com/milvus-io/milvus-operator/apis/milvus.io/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Test(t *testing.T) {
	s := file.NewStorage("")
	cfg, err := s.GetKubeConfigByCluster("")
	assert.NoError(t, err)
	c, err := NewClient("", cfg)
	assert.NoError(t, err)

	mcList := &v1alpha1.MilvusClusterList{}
	err = c.rawCli.List(context.TODO(), mcList, &client.ListOptions{Namespace: "default"})
	assert.NoError(t, err)

	deployList := &appsv1.DeploymentList{}
	err = c.rawCli.List(context.TODO(), deployList, &client.ListOptions{Namespace: "default"})
	assert.NoError(t, err)

	svcList := &corev1.ServiceList{}
	err = c.rawCli.List(context.TODO(), svcList, &client.ListOptions{Namespace: "default"})
	assert.NoError(t, err)
}
