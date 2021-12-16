package k8s

import (
	"context"
	"sync"

	"github.com/haorenfsa/milvus-ops/service"
	"github.com/haorenfsa/milvus-ops/storage"
	"github.com/pkg/errors"
)

type K8sClient = service.K8sClient

type K8sClientGetter struct {
	storage storage.Kubeconfig
	cache   map[string]K8sClient
	lock    sync.Mutex
}

// NewK8sClientGetter implements K8sClientGetter interface by file
func NewK8sClientGetter(storage storage.Kubeconfig) *K8sClientGetter {
	return &K8sClientGetter{storage: storage, cache: map[string]K8sClient{}}
}

func (k *K8sClientGetter) GetClientByCluster(ctx context.Context, cluster string) (K8sClient, error) {
	k.lock.Lock()
	defer k.lock.Unlock()
	if client, ok := k.cache[cluster]; ok {
		return client, nil
	}
	kubeconfig, err := k.storage.GetKubeConfigByCluster(cluster)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kubeconfig")
	}
	cli, err := NewClient(cluster, kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create k8s client")
	}
	k.cache[cluster] = cli
	return cli, nil
}