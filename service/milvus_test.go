package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/haorenfsa/milvus-ops/model"
	"github.com/stretchr/testify/assert"
)

func TestMilvusServiceListAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockK8s := NewMockK8sClientGetter(ctrl)
	mockHelm := NewMockHelmClientForMilvus(ctrl)
	mockK8sCli := NewMockK8sClient(ctrl)

	milvusService := NewMilvusService(mockK8s, mockHelm)

	ctx := context.TODO()

	cluster := "cluster"

	testErr := errors.New("test")

	// get client err
	mockK8s.EXPECT().GetClientByCluster(gomock.Any(), cluster).
		Return(nil, testErr)
	_, err := milvusService.ListAll(ctx, cluster)
	assert.Error(t, err)

	// list ns err
	mockK8s.EXPECT().GetClientByCluster(gomock.Any(), cluster).
		Return(mockK8sCli, nil)
	mockK8sCli.EXPECT().ListNamespaces(gomock.Any()).Return(nil, testErr)
	_, err = milvusService.ListAll(ctx, cluster)
	assert.Error(t, err)

	nss := []string{"ns1", "ns2"}
	// list helm err
	mockK8s.EXPECT().GetClientByCluster(gomock.Any(), cluster).
		Return(mockK8sCli, nil)
	mockK8sCli.EXPECT().ListNamespaces(gomock.Any()).Return(nss, nil)
	mockHelm.EXPECT().ListMilvus(gomock.Any(), mockK8sCli, nss[0]).Return(nil, testErr)
	_, err = milvusService.ListAll(ctx, cluster)
	assert.Error(t, err)

	// list crd err
	milvusInHelm1 := []*Milvus{
		{Name: "milvus1", Namespace: nss[0]},
		{Name: "milvus2", Namespace: nss[0]},
	}
	milvusInHelm2 := []*Milvus{}
	mockK8s.EXPECT().GetClientByCluster(gomock.Any(), cluster).
		Return(mockK8sCli, nil)
	mockK8sCli.EXPECT().ListNamespaces(gomock.Any()).Return(nss, nil)
	mockHelm.EXPECT().ListMilvus(gomock.Any(), mockK8sCli, nss[0]).Return(milvusInHelm1, nil)
	mockHelm.EXPECT().ListMilvus(gomock.Any(), mockK8sCli, nss[1]).Return(milvusInHelm2, nil)
	mockK8sCli.EXPECT().ListMilvusCluster(gomock.Any(), "").Return(nil, testErr)
	_, err = milvusService.ListAll(ctx, cluster)
	assert.Error(t, err)

	// all ok
	milvusCrd := []*Milvus{
		{Name: "milvus3", Namespace: nss[1]},
	}
	mockK8s.EXPECT().GetClientByCluster(gomock.Any(), cluster).
		Return(mockK8sCli, nil)
	mockK8sCli.EXPECT().ListNamespaces(gomock.Any()).Return(nss, nil)
	mockHelm.EXPECT().ListMilvus(gomock.Any(), mockK8sCli, nss[0]).Return(milvusInHelm1, nil)
	mockHelm.EXPECT().ListMilvus(gomock.Any(), mockK8sCli, nss[1]).Return(milvusInHelm2, nil)
	mockK8sCli.EXPECT().ListMilvusCluster(gomock.Any(), "").Return(milvusCrd, nil)
	ret, err := milvusService.ListAll(ctx, cluster)
	assert.NoError(t, err)
	assert.Len(t, ret, 3)
}
