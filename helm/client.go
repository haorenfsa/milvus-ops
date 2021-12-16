package helm

import (
	"context"
	"strings"
	"time"

	. "github.com/haorenfsa/milvus-ops/model"
	"github.com/haorenfsa/milvus-ops/service"
	"github.com/pkg/errors"
	"github.com/tevino/log"
	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"

	"github.com/patrickmn/go-cache"
)

type ChartRequest struct {
	ReleaseName string
	Namespace   string
	Chart       string
	Values      map[string]interface{}
}

func NeedUpdate(status release.Status) bool {
	return status == release.StatusFailed ||
		status == release.StatusUnknown ||
		status == release.StatusUninstalled
}

// LocalClient is the local implementation of the Client interface.
type LocalClient struct{}

func (d *LocalClient) GetStatus(cfg *action.Configuration, releaseName string) (release.Status, error) {
	client := action.NewStatus(cfg)
	rel, err := client.Run(releaseName)
	if err != nil {
		return release.StatusUnknown, err
	}

	return rel.Info.Status, nil
}

func (d *LocalClient) GetValues(cfg *action.Configuration, releaseName string) (map[string]interface{}, error) {
	client := action.NewGetValues(cfg)
	vals, err := client.Run(releaseName)
	if err != nil {
		return nil, err
	}

	if vals == nil {
		return map[string]interface{}{}, nil
	}

	return vals, nil
}

func (d *LocalClient) ReleaseExist(cfg *action.Configuration, releaseName string) (bool, error) {
	histClient := action.NewHistory(cfg)
	histClient.Max = 1
	_, err := histClient.Run(releaseName)
	if err == driver.ErrReleaseNotFound {
		return false, nil
	}

	return err == nil, err
}

func (d *LocalClient) Upgrade(cfg *action.Configuration, request ChartRequest) error {
	exist, err := d.ReleaseExist(cfg, request.ReleaseName)
	if err != nil {
		return err
	}
	if !exist {
		return d.Install(cfg, request)
	}

	return d.Update(cfg, request)
}

func (d *LocalClient) Update(cfg *action.Configuration, request ChartRequest) error {
	client := action.NewUpgrade(cfg)
	client.Namespace = request.Namespace
	chartRequested, err := loader.Load(request.Chart)
	if err != nil {
		return err
	}
	if len(request.Values) == 0 {
		client.ResetValues = true
	}

	_, err = client.Run(request.ReleaseName, chartRequested, request.Values)
	return err
}

func (d *LocalClient) Install(cfg *action.Configuration, request ChartRequest) error {
	client := action.NewInstall(cfg)
	client.ReleaseName = request.ReleaseName
	client.Namespace = request.Namespace
	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}

	chartRequested, err := loader.Load(request.Chart)
	if err != nil {
		return err
	}

	_, err = client.Run(chartRequested, request.Values)
	return err
}

func (d *LocalClient) Uninstall(cfg *action.Configuration, releaseName string) error {
	_, err := cfg.Releases.History(releaseName)
	if errors.Is(err, driver.ErrReleaseNotFound) {
		return nil
	}

	client := action.NewUninstall(cfg)
	client.DisableHooks = true
	_, err = client.Run(releaseName)
	if err != nil {
		return err
	}

	return nil
}

// Client for helm
type Client struct {
	*LocalClient
	cache *cache.Cache
}

func NewClients(clis []service.K8sClient) *Client {
	ret := &Client{
		cache:       cache.New(time.Second*180, time.Second*5),
		LocalClient: &LocalClient{},
	}
	for _, iter := range clis {
		cli := iter
		go ret.SyncRelease(context.Background(), cli)
	}
	return ret
}

func (c *Client) ListMilvus(ctx context.Context, client service.K8sClient, namespace string) ([]*Milvus, error) {
	cacheKey := strings.Join([]string{client.ClusterName(), namespace, "milvus"}, ".")
	res, found := c.cache.Get(cacheKey)
	if !found {
		return nil, errors.New("server still loading data")
	}
	releases := res.([]*release.Release)
	ret := releasesToMilvus(releases)
	return ret, nil
}

func (c *Client) list(ctx context.Context, client service.K8sClient, namespace string) ([]*release.Release, error) {
	cfg := new(action.Configuration)
	cfg.Init(client.RESTClientGetter(), namespace, "", log.Infof)
	listAction := action.NewList(cfg)
	releases, err := listAction.Run()
	if err != nil {
		return nil, errors.Wrap(err, "listAction run failed")
	}
	return releases, nil
}

func (c *Client) SyncRelease(ctx context.Context, client service.K8sClient) {
	ticker := time.NewTicker(time.Minute * 2)
	for {
		nss, err := client.ListNamespaces(ctx)
		if err != nil {
			log.Errorf("list namespace failed: %v", err)
			continue
		}
		eg := errgroup.Group{}
		log.Info("start sync release")
		start := time.Now()
		for _, iter := range nss {
			ns := iter
			eg.Go(func() error {
				releases, err := c.list(ctx, client, ns)
				if err != nil {
					return err
				}
				cacheKey := strings.Join([]string{client.ClusterName(), ns, "milvus"}, ".")
				c.cache.Set(cacheKey, releases, time.Hour)
				return nil
			})
		}
		err = eg.Wait()
		log.Info("sync release done, cost: ", time.Since(start))
		if err != nil {
			log.Errorf("sync release failed: %v", err)
		} else {
			log.Info("sync release success")
		}
		select {
		case <-ticker.C:

		case <-ctx.Done():
			return
		}

	}
}

func releasesToMilvus(releases []*release.Release) []*Milvus {
	var milvuses []*Milvus
	for _, release := range releases {
		if !strings.HasPrefix(release.Chart.Metadata.Name, "milvus") {
			continue
		}
		milvuses = append(milvuses, releaseToMilvus(release))
	}
	return milvuses
}

func releaseToMilvus(release *release.Release) *Milvus {
	return &Milvus{
		Name:      release.Name,
		Namespace: release.Namespace,
		Status:    release.Info.Status.String(), //TODO: sync real status
		Version:   release.Chart.Metadata.AppVersion,
		ManagedBy: "helm",
		// Type:
	}
}
