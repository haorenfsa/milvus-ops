package file

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Storage by db
type Storage struct {
	rootPath string
}

func NewStorage(rootPath string) *Storage {
	if rootPath == "" {
		rootPath = os.Getenv("HOME") + "/.kube"
	}
	return &Storage{rootPath}
}

func (s Storage) ListClusters() ([]string, error) {
	files, err := ioutil.ReadDir(s.rootPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list clusters")
	}

	var clusters []string
	for _, file := range files {
		if !file.IsDir() {
			if file.Name() == "config" {
				clusters = append(clusters, "default")
			} else {
				cluster := strings.TrimSuffix(file.Name(), ".yaml")
				clusters = append(clusters, cluster)
			}
		}

	}
	return clusters, nil
}

func (s Storage) GetKubeConfigByCluster(cluster string) ([]byte, error) {
	var filePath string
	switch cluster {
	case "", "default":
		filePath = s.rootPath + "/config"
	default:
		filePath = s.rootPath + "/" + cluster + ".yaml"
	}
	data, err := ioutil.ReadFile(filePath)
	return data, errors.Wrapf(err, "failed to read kubeconfig from %s", filePath)
}

type K8s struct {
	ID         int64
	Name       string
	Kubeconfig []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
