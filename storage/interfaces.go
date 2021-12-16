package storage

type Kubeconfig interface {
	GetKubeConfigByCluster(cluster string) ([]byte, error)
}
