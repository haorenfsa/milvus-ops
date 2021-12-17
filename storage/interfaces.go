package storage

type Kubeconfig interface {
	ListClusters() ([]string, error)
	GetKubeConfigByCluster(cluster string) ([]byte, error)
}
