package service

import (
	"context"
	"io"

	. "github.com/haorenfsa/milvus-ops/model"
	"github.com/maoqide/kubeutil/pkg/terminal"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

//go:generate mockgen -source=dependencies.go -destination=dependencies_mock.go -package=service

type K8sClientGetter interface {
	GetClientByCluster(ctx context.Context, cluster string) (K8sClient, error)
	ListClients(ctx context.Context) ([]K8sClient, error)
	ListClusters(ctx context.Context) ([]string, error)
}

type K8sClient interface {
	ClusterName() string
	Shell(ctx context.Context, hdl terminal.PtyHandler, loc MilvusLocateOption) error
	RESTClientGetter() genericclioptions.RESTClientGetter

	ListNamespaces(ctx context.Context) ([]string, error)
	// ListMilvusCluster return by namespace, empty namespace means all namespaces
	ListMilvusCluster(ctx context.Context, namespace string) ([]*Milvus, error)
	// ListPods return pods
	ListPods(ctx context.Context, opt MilvusLocateOption) ([]string, error)
	// ListPods return pods
	ListPodsDetail(ctx context.Context, opt MilvusLocateOption) ([]corev1.Pod, error)
	// ListPodsByLabel return logstream by options
	Logs(ctx context.Context, ptyHandler terminal.PtyHandler, opt MilvusLocateOption) error

	DownloadLog(ctx context.Context, opt MilvusLocateOption) (io.ReadCloser, error)
}

type HelmClientForMilvus interface {
	ListMilvus(ctx context.Context, client K8sClient, namespace string) ([]*Milvus, error)
}
