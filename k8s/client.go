package k8s

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	. "github.com/haorenfsa/milvus-ops/model"
	"github.com/haorenfsa/milvus-ops/util"
	"github.com/maoqide/kubeutil/pkg/terminal"
	"github.com/pkg/errors"
	"github.com/tevino/log"
	"golang.org/x/sync/errgroup"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/milvus-io/milvus-operator/apis/milvus.io/v1alpha1"
)

type Client struct {
	name             string
	rawCli           client.Client // with cache, feel free to use it
	restCfg          rest.Config
	wrappedCliGetter *WrappedRestClientGetter
}

func NewClient(cluster string, kubeconfig []byte) (*Client, error) {
	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client config")
	}
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rest config")
	}
	scheme := runtime.NewScheme()
	err = v1alpha1.SchemeBuilder.AddToScheme(scheme)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add v1alpha1 scheme")
	}
	err = appsv1.AddToScheme(scheme)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add appsv1 scheme")
	}
	err = corev1.AddToScheme(scheme)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add corev1 scheme")
	}
	rawCli, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client by rest config")
	}
	dis, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create discovery client")
	}
	cachedDis := memory.NewMemCacheClient(dis)
	mapper := rawCli.RESTMapper()
	cliGetter := NewWrappedRestClientGetter(config, mapper, cachedDis, clientConfig)
	return &Client{cluster, rawCli, *config, cliGetter}, nil
}

func (c *Client) ClusterName() string {
	return c.name
}

func (c *Client) ListNamespaces(ctx context.Context) ([]string, error) {
	nsList := &corev1.NamespaceList{}
	err := c.rawCli.List(ctx, nsList)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list namespaces")
	}
	var ret = []string{}
	for _, v := range nsList.Items {
		ret = append(ret, v.Name)
	}
	return ret, nil
}

func (c *Client) ListMilvusCluster(ctx context.Context, namespace string) ([]*Milvus, error) {
	mcList := &v1alpha1.MilvusClusterList{}
	err := c.rawCli.List(ctx, mcList, client.InNamespace(namespace))
	if meta.IsNoMatchError(err) {
		return []*Milvus{}, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to list milvus clusters")
	}
	return convertMilvusList(mcList), nil
}

func (c *Client) DownloadLog(ctx context.Context, opt MilvusLocateOption) (io.ReadCloser, error) {
	pods, err := c.ListPods(ctx, opt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list pods")
	}

	if len(pods) > 1 {
		return nil, errors.New("not support multiple pods")
	}
	pod := pods[0]

	cfg := c.restCfg
	cfg.APIPath = "api"
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	cfg.NegotiatedSerializer = scheme.Codecs

	// restCli, err := rest.RESTClientFor(&cfg)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to create rest client")
	// }

	cliSet, err := kubernetes.NewForConfig(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %v", err)
	}

	req := cliSet.CoreV1().
		Pods(opt.Namespace).
		GetLogs(pod, &corev1.PodLogOptions{
			Container: opt.Container,
			Follow:    false,
			Previous:  false,
		})

	return req.Stream(ctx)
}

func (c *Client) Logs(ctx context.Context, ptyHandler terminal.PtyHandler, opt MilvusLocateOption) error {

	pods, err := c.ListPods(ctx, opt)
	if err != nil {
		return errors.Wrap(err, "failed to list pods")
	}

	cfg := c.restCfg
	cfg.APIPath = "api"
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	cfg.NegotiatedSerializer = scheme.Codecs

	// restCli, err := rest.RESTClientFor(&cfg)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to create rest client")
	// }

	cliSet, err := kubernetes.NewForConfig(&cfg)
	if err != nil {
		return fmt.Errorf("unable to create a client: %v", err)
	}

	logWriter := NewConcurrentWriter(ptyHandler.Stdout())
	tailLine := int64(300)
	var logOnePod = func(podName string) error {
		req := cliSet.CoreV1().
			Pods(opt.Namespace).
			GetLogs(podName, &corev1.PodLogOptions{
				Container: opt.Container,
				Follow:    true,
				TailLines: &tailLine,
			})

		stream, err := req.Stream(ctx)
		if err != nil {
			return err
		}
		defer stream.Close()

		_, err = io.Copy(logWriter, stream)

		if err == io.EOF {
			return nil
		}

		// executor, err := remotecommand.NewSPDYExecutor(&cfg, "GET", req.URL())
		// if err != nil {
		// 	return errors.Wrap(err, "failed to create executor")
		// }

		// // Stream
		// err = executor.Stream(remotecommand.StreamOptions{
		// 	Stdin:             ptyHandler.Stdin(),
		// 	Stdout:            ptyHandler.Stdout(),
		// 	Stderr:            ptyHandler.Stderr(),
		// 	TerminalSizeQueue: ptyHandler,
		// 	Tty:               false,
		// })
		return errors.Wrap(err, "failed to stream")
	}

	eg := errgroup.Group{}
	for _, iter := range pods {
		pod := iter
		eg.Go(func() error {
			log.Info("logging ", pod)
			return errors.Wrap(logOnePod(pod), "failed to log one pod")
		})
	}
	return eg.Wait()
}

func (c *Client) ListPods(ctx context.Context, opt MilvusLocateOption) ([]string, error) {
	if opt.Pod != "" {
		return []string{opt.Pod}, nil
	}
	podList := &corev1.PodList{}
	labels := client.MatchingLabels{
		"app.kubernetes.io/instance": opt.Milvus,
		"app.kubernetes.io/name":     "milvus",
	}
	if opt.Component != "" {
		label := util.GetComponentLabelByManager(opt.ManagedBy)
		labels[label] = opt.Component
	}
	err := c.rawCli.List(ctx, podList,
		client.InNamespace(opt.Namespace),
		labels,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list pods")
	}
	var ret = []string{}
	for _, v := range podList.Items {
		ret = append(ret, v.Name)
	}
	return ret, nil
}

func (c *Client) ListPodsDetail(ctx context.Context, opt MilvusLocateOption) ([]corev1.Pod, error) {
	podList := &corev1.PodList{}
	labels := client.MatchingLabels{
		"app.kubernetes.io/instance": opt.Milvus,
		"app.kubernetes.io/name":     "milvus",
	}
	if opt.Component != "" {
		label := util.GetComponentLabelByManager(opt.ManagedBy)
		labels[label] = opt.Component
	}
	err := c.rawCli.List(ctx, podList,
		client.InNamespace(opt.Namespace),
		labels,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list pods")
	}
	return podList.Items, nil
}

func (c *Client) Shell(ctx context.Context, ptyHandler terminal.PtyHandler, loc MilvusLocateOption) error {
	// cmd := []string{"/bin/sh"}
	cmd := []string{"/bin/sh", "-c", "/bin/bash || /bin/sh"}

	// pod := corev1.Pod{}
	// err := c.rawCli.Get(ctx, types.NamespacedName{Namespace: loc.Namespace, Name: loc.Pod}, &pod)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to found pod")
	// }

	if loc.Pod == "" {
		pods, err := c.ListPods(ctx, loc)
		if err != nil {
			return errors.Wrap(err, "failed to list pods")
		}
		if len(pods) < 1 {
			return errors.New("no pod found")
		}
		loc.Pod = pods[0]
		if len(pods) > 1 {
			for _, pod := range pods {
				if strings.Contains(pod, "proxy") {
					loc.Pod = pod
					break
				}
			}
		}
	}

	cfg := c.restCfg
	cfg.APIPath = "api"
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	cfg.NegotiatedSerializer = scheme.Codecs
	restCli, err := rest.RESTClientFor(&cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create rest client")
	}
	req := restCli.Post().
		Namespace(loc.Namespace).
		Resource("pods").
		Name(loc.Pod).
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		Container: loc.Container,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(&cfg, "POST", req.URL())
	if err != nil {
		return errors.Wrap(err, "failed to create executor")
	}

	// Stream
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             ptyHandler.Stdin(),
		Stdout:            ptyHandler.Stdout(),
		Stderr:            ptyHandler.Stderr(),
		TerminalSizeQueue: ptyHandler,
		Tty:               true,
	})
	return errors.Wrap(err, "failed to stream")
}

func convertMilvusList(mcList *v1alpha1.MilvusClusterList) []*Milvus {
	var ret = []*Milvus{}
	for _, v := range mcList.Items {
		ret = append(ret, convertMilvus(&v))
	}
	return ret
}

func convertMilvus(mc *v1alpha1.MilvusCluster) *Milvus {
	return &Milvus{
		Namespace: mc.Namespace,
		Name:      mc.Name,
		Status:    string(mc.Status.Status),
		Version:   getImageTag(mc.Spec.Com.Image),
		ManagedBy: "operator",
	}
}

func getImageTag(image string) string {
	splited := strings.Split(image, ":")
	if len(splited) < 2 {
		return "latest"
	}
	return splited[1]
}

// ValidateContainer validate container.
func ValidateContainer(pod *corev1.Pod, containerName string) (bool, error) {
	if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
		return false, fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
	}
	for _, c := range pod.Spec.Containers {
		if containerName == c.Name {
			return true, nil
		}
	}
	return false, fmt.Errorf("pod has no container '%s'", containerName)
}

type WrappedRestClientGetter struct {
	restConfig               *rest.Config
	mapper                   meta.RESTMapper
	cachedDiscoveryInterface discovery.CachedDiscoveryInterface
	clientconfig             clientcmd.ClientConfig
}

func NewWrappedRestClientGetter(restConfig *rest.Config, mapper meta.RESTMapper, cachedDiscoveryInterface discovery.CachedDiscoveryInterface, clientconfig clientcmd.ClientConfig) *WrappedRestClientGetter {
	return &WrappedRestClientGetter{
		restConfig:               restConfig,
		mapper:                   mapper,
		cachedDiscoveryInterface: cachedDiscoveryInterface,
		clientconfig:             clientconfig,
	}
}

func (w *WrappedRestClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	dis, err := discovery.NewDiscoveryClientForConfig(w.restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create discovery client")
	}
	return memory.NewMemCacheClient(dis), nil
}

func (w *WrappedRestClientGetter) ToRESTConfig() (*rest.Config, error) {
	return w.restConfig, nil
}

func (w *WrappedRestClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	return meta.NewDefaultRESTMapper([]schema.GroupVersion{*w.restConfig.GroupVersion}), nil
}

func (w *WrappedRestClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return w.clientconfig
}

func (c Client) RESTClientGetter() genericclioptions.RESTClientGetter {
	return c.wrappedCliGetter
}

// ConcurrentWriter helps a non concurrent writer to write concurrently with a lock
type ConcurrentWriter struct {
	writer    io.Writer
	writeLock sync.Mutex
}

func NewConcurrentWriter(w io.Writer) *ConcurrentWriter {
	return &ConcurrentWriter{
		writer: w,
	}
}

func (c *ConcurrentWriter) Write(p []byte) (int, error) {
	c.writeLock.Lock()
	ret, err := c.writer.Write(p)
	c.writeLock.Unlock()
	return ret, err
}
