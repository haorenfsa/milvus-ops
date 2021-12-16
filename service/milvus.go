package service

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/haorenfsa/milvus-ops/model"
	"github.com/haorenfsa/milvus-ops/util"
	websockTerm "github.com/maoqide/kubeutil/pkg/terminal/websocket"
	"github.com/pkg/errors"
	"github.com/tevino/log"
)

type MilvusService struct {
	clusters K8sClientGetter
	helm     HelmClientForMilvus
}

func NewMilvusService(clusters K8sClientGetter, helm HelmClientForMilvus) *MilvusService {
	return &MilvusService{
		clusters: clusters,
		helm:     helm,
	}
}

const allNamespaces = ""

var ErrServerConfig = errors.New("server config error")
var ErrK8sAPI = errors.New("call k8s api error")
var ErrHelm = errors.New("call helm error")

func wrapErr(errType error, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	msg := fmt.Sprintf(format, args...)
	err = errors.Wrapf(errType, "%s: %s", msg, err)
	return err
}

func (m *MilvusService) ListNamespaces(ctx context.Context, cluster string) ([]string, error) {
	k8sCli, err := m.clusters.GetClientByCluster(ctx, cluster)
	if err != nil {
		return nil, wrapErr(ErrServerConfig, err, "get k8s client for [cluster=%s]", cluster)
	}
	nss, err := k8sCli.ListNamespaces(ctx)
	return nss, wrapErr(ErrK8sAPI, err, "list namespaces for [cluster=%s]", cluster)
}

func (m *MilvusService) DownloadLog(ctx context.Context, opt MilvusLocateOption) (io.ReadCloser, error) {
	k8sCli, err := m.clusters.GetClientByCluster(ctx, opt.Cluster)
	if err != nil {
		return nil, wrapErr(ErrServerConfig, err, "get k8s client for [cluster=%s]", opt.Cluster)
	}
	return k8sCli.DownloadLog(ctx, opt)
}

func (m *MilvusService) Logs(ctx context.Context, wsConn *websocket.Conn, opt MilvusLocateOption) error {
	k8sCli, err := m.clusters.GetClientByCluster(ctx, opt.Cluster)
	if err != nil {
		return wrapErr(ErrServerConfig, err, "get k8s client for [cluster=%s]", opt.Cluster)
	}
	err = k8sCli.Logs(ctx, websockTerm.NewTerminalSessionWs(wsConn), opt)
	return wrapErr(ErrK8sAPI, err, "logs for [cluster=%s,ns=%s,milvus=%s]", opt.Cluster, opt.Namespace, opt.Milvus)
}

func (m *MilvusService) ListPods(ctx context.Context, opt MilvusLocateOption) (*ClassifiedPods, error) {
	k8sCli, err := m.clusters.GetClientByCluster(ctx, opt.Cluster)
	if err != nil {
		return nil, wrapErr(ErrServerConfig, err, "get k8s client for [cluster=%s]", opt.Cluster)
	}
	pods, err := k8sCli.ListPodsDetail(ctx, opt)
	if err != nil {
		wrapErr(ErrK8sAPI, err, "list pods for [cluster=%s,ns=%s,milvus=%s]", opt.Cluster, opt.Namespace, opt.Milvus)
	}
	ret := NewClassifiedPods()
	label := util.GetComponentLabelByManager(opt.ManagedBy)
	for _, pod := range pods {
		component := pod.Labels[label]
		comPods, existed := ret.ComponentPods[component]
		if !existed {
			comPods = []string{}
			ret.Components = append(ret.Components, component)
		}
		comPods = append(comPods, pod.Name)
		ret.ComponentPods[component] = comPods
	}
	return ret, nil
}

func (m *MilvusService) GetComponents(ctx context.Context, opt MilvusLocateOption) ([]string, error) {
	k8sCli, err := m.clusters.GetClientByCluster(ctx, opt.Cluster)
	if err != nil {
		return nil, wrapErr(ErrServerConfig, err, "get k8s client for [cluster=%s]", opt.Cluster)
	}
	pods, err := k8sCli.ListPods(ctx, opt)
	return pods, wrapErr(ErrK8sAPI, err, "list pods for [cluster=%s,ns=%s,milvus=%s]", opt.Cluster, opt.Namespace, opt.Milvus)
}

func (m *MilvusService) ListAll(ctx context.Context, cluster string) ([]*Milvus, error) {
	k8sCli, err := m.clusters.GetClientByCluster(ctx, cluster)
	if err != nil {
		return nil, wrapErr(ErrServerConfig, err, "get k8s client for [cluster=%s]", cluster)
	}

	nss, err := k8sCli.ListNamespaces(ctx)
	if err != nil {
		return nil, wrapErr(ErrK8sAPI, err, "list namespaces for [cluster=%s]", cluster)
	}

	var milvusList []*Milvus
	var errs []error
	// list by helm
	start := time.Now()
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, iter := range nss {
		ns := iter
		wg.Add(1)
		go func() {
			defer wg.Done()
			tmpList, err := m.helm.ListMilvus(ctx, k8sCli, ns)
			lock.Lock()
			defer lock.Unlock()
			if err != nil {
				err = wrapErr(ErrHelm, err, "helm list milvus for [cluster=%s,ns=%s]", cluster, ns)
				errs = append(errs, err)
				return
			}
			milvusList = append(milvusList, tmpList...)
		}()
	}
	wg.Wait()
	if len(errs) > 0 {
		return nil, errors.Wrapf(ErrHelm, "list milvus error: %v", errs)
	}
	log.Info("list milvus by helm", "time", time.Since(start))

	// list by crd
	start = time.Now()
	tmpList, err := k8sCli.ListMilvusCluster(ctx, allNamespaces)
	if err != nil {
		return nil, wrapErr(ErrK8sAPI, err, "list milvus for [cluster=%s,ns=%s]", cluster, allNamespaces)
	}
	log.Info("list milvus by crd", "time", time.Since(start))
	milvusList = append(milvusList, tmpList...)
	sortMilvus := sortMilvus(milvusList)
	sort.Sort(&sortMilvus)
	return milvusList, nil
}

type sortMilvus []*Milvus

func (s *sortMilvus) Len() int {
	return len(*s)
}

func (s *sortMilvus) Less(i, j int) bool {
	return (*s)[i].Name < (*s)[j].Name
}

func (s *sortMilvus) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func (m *MilvusService) Shell(ctx context.Context, wsConn *websocket.Conn, loc MilvusLocateOption) error {
	k8sCli, err := m.clusters.GetClientByCluster(ctx, loc.Cluster)
	if err != nil {
		return wrapErr(ErrServerConfig, err, "get k8s client for [cluster=%s]", loc.Cluster)
	}
	err = k8sCli.Shell(ctx, websockTerm.NewTerminalSessionWs(wsConn), loc)
	return wrapErr(ErrK8sAPI, err, "shell for [cluster=%s,ns=%s,milvus=%s,component=%s]", loc.Cluster, loc.Namespace, loc.Milvus, loc.Component)
}
