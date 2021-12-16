// Code generated by MockGen. DO NOT EDIT.
// Source: dependencies.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	. "github.com/haorenfsa/milvus-ops/model"
)

// MockK8sClientGetter is a mock of K8sClientGetter interface.
type MockK8sClientGetter struct {
	ctrl     *gomock.Controller
	recorder *MockK8sClientGetterMockRecorder
}

// MockK8sClientGetterMockRecorder is the mock recorder for MockK8sClientGetter.
type MockK8sClientGetterMockRecorder struct {
	mock *MockK8sClientGetter
}

// NewMockK8sClientGetter creates a new mock instance.
func NewMockK8sClientGetter(ctrl *gomock.Controller) *MockK8sClientGetter {
	mock := &MockK8sClientGetter{ctrl: ctrl}
	mock.recorder = &MockK8sClientGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockK8sClientGetter) EXPECT() *MockK8sClientGetterMockRecorder {
	return m.recorder
}

// GetClientByCluster mocks base method.
func (m *MockK8sClientGetter) GetClientByCluster(ctx context.Context, cluster string) (K8sClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClientByCluster", ctx, cluster)
	ret0, _ := ret[0].(K8sClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClientByCluster indicates an expected call of GetClientByCluster.
func (mr *MockK8sClientGetterMockRecorder) GetClientByCluster(ctx, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClientByCluster", reflect.TypeOf((*MockK8sClientGetter)(nil).GetClientByCluster), ctx, cluster)
}

// MockK8sClient is a mock of K8sClient interface.
type MockK8sClient struct {
	ctrl     *gomock.Controller
	recorder *MockK8sClientMockRecorder
}

// MockK8sClientMockRecorder is the mock recorder for MockK8sClient.
type MockK8sClientMockRecorder struct {
	mock *MockK8sClient
}

// NewMockK8sClient creates a new mock instance.
func NewMockK8sClient(ctrl *gomock.Controller) *MockK8sClient {
	mock := &MockK8sClient{ctrl: ctrl}
	mock.recorder = &MockK8sClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockK8sClient) EXPECT() *MockK8sClientMockRecorder {
	return m.recorder
}

// ListMilvusCluster mocks base method.
func (m *MockK8sClient) ListMilvusCluster(ctx context.Context, namespace string) ([]*Milvus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMilvusCluster", ctx, namespace)
	ret0, _ := ret[0].([]*Milvus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMilvusCluster indicates an expected call of ListMilvusCluster.
func (mr *MockK8sClientMockRecorder) ListMilvusCluster(ctx, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMilvusCluster", reflect.TypeOf((*MockK8sClient)(nil).ListMilvusCluster), ctx, namespace)
}

// ListNamespaces mocks base method.
func (m *MockK8sClient) ListNamespaces(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNamespaces", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNamespaces indicates an expected call of ListNamespaces.
func (mr *MockK8sClientMockRecorder) ListNamespaces(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNamespaces", reflect.TypeOf((*MockK8sClient)(nil).ListNamespaces), ctx)
}

// MockHelmClientForMilvus is a mock of HelmClientForMilvus interface.
type MockHelmClientForMilvus struct {
	ctrl     *gomock.Controller
	recorder *MockHelmClientForMilvusMockRecorder
}

// MockHelmClientForMilvusMockRecorder is the mock recorder for MockHelmClientForMilvus.
type MockHelmClientForMilvusMockRecorder struct {
	mock *MockHelmClientForMilvus
}

// NewMockHelmClientForMilvus creates a new mock instance.
func NewMockHelmClientForMilvus(ctrl *gomock.Controller) *MockHelmClientForMilvus {
	mock := &MockHelmClientForMilvus{ctrl: ctrl}
	mock.recorder = &MockHelmClientForMilvusMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHelmClientForMilvus) EXPECT() *MockHelmClientForMilvusMockRecorder {
	return m.recorder
}

// ListMilvus mocks base method.
func (m *MockHelmClientForMilvus) ListMilvus(ctx context.Context, client K8sClient, namespace string) ([]*Milvus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMilvus", ctx, client, namespace)
	ret0, _ := ret[0].([]*Milvus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMilvus indicates an expected call of ListMilvus.
func (mr *MockHelmClientForMilvusMockRecorder) ListMilvus(ctx, client, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMilvus", reflect.TypeOf((*MockHelmClientForMilvus)(nil).ListMilvus), ctx, client, namespace)
}
