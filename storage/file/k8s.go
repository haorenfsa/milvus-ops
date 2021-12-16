package file

import (
	"io/ioutil"
	"os"
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

func (s Storage) GetKubeConfigByCluster(cluster string) ([]byte, error) {
	var filePath string
	switch cluster {
	case "":
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
