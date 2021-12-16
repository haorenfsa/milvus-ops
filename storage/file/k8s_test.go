package file

import (
	"os"
	"testing"

	"github.com/haorenfsa/milvus-ops/util"
	"github.com/stretchr/testify/assert"
)

func TestStorage_NewStorage(t *testing.T) {
	s := NewStorage("")
	assert.Equal(t, os.Getenv("HOME")+"/.kube", s.rootPath)
	s = NewStorage("./")
	assert.Equal(t, "./", s.rootPath)
}

func TestStorage_GetKubeConfigByCluster(t *testing.T) {
	s := NewStorage(util.GetGitRepoRootDir() + "/test/kubeconfig/")
	ret, err := s.GetKubeConfigByCluster("")
	assert.NoError(t, err)
	assert.Equal(t, []byte("name: default"), ret)

	ret, err = s.GetKubeConfigByCluster("cluster")
	assert.NoError(t, err)
	assert.Equal(t, []byte("name: cluster"), ret)

	_, err = s.GetKubeConfigByCluster("not-exist")
	assert.Error(t, err)
}
