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

func (k *K8sClientGetter) ListClients(ctx context.Context) ([]K8sClient, error) {
	clusters, err := k.storage.ListClusters()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list clusters")
	}
	var clients []K8sClient
	for _, cluster := range clusters {
		cli, err := k.GetClientByCluster(ctx, cluster)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get k8s client")
		}
		clients = append(clients, cli)
	}
	return clients, nil
}

func (k *K8sClientGetter) ListClusters(ctx context.Context) ([]string, error) {
	var ret = []string{}
	for name := range k.cache {
		ret = append(ret, name)
	}
	return ret, nil
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
		return nil, errors.Wrapf(err, "failed to create k8s client for cluster %s", cluster)
	}
	k.cache[cluster] = cli
	return cli, nil
}
