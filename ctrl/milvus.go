package ctrl

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/haorenfsa/milvus-ops/model"
	"github.com/haorenfsa/milvus-ops/server"
	"github.com/haorenfsa/milvus-ops/service"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/tevino/log"
)

type MilvusController struct {
	service *service.MilvusService
}

func NewMilvusController(s *service.MilvusService) *MilvusController {
	return &MilvusController{
		service: s,
	}
}

// Register registers request handler
func (c *MilvusController) Register(root gin.IRouter) {
	root.GET("clusters", server.WrapHandler(c.handleListClusters))
	g := root.Group("/clusters/:cluster")
	g.GET("/namespaces", server.WrapHandler(c.handleGetNamespaces))
	g.GET("/milvus", server.WrapHandler(c.handleGet))
	g.GET("/milvus/:namespace/:milvus/pods", server.WrapHandler(c.handleGetPods))
	g.GET("/milvus/:namespace/:milvus/shell", c.handleShell)
	g.GET("/milvus/:namespace/:milvus/logs", c.handleGetLogs)
	g.GET("/milvus/:namespace/:milvus/files/log", c.handleDownloadLog)
}

func (c *MilvusController) handleListClusters(ctx *gin.Context) (interface{}, error) {
	return c.service.ListClusters(ctx)
}

func (c *MilvusController) handleGetNamespaces(ctx *gin.Context) (interface{}, error) {
	opt := model.GetMilvusLocateOption(ctx)
	return c.service.ListNamespaces(ctx, opt.Cluster)
}

func (c *MilvusController) handleGet(ctx *gin.Context) (interface{}, error) {
	opt := model.GetMilvusLocateOption(ctx)
	cacheKey := strings.Join([]string{opt.Cluster, opt.Namespace, "milvus"}, ".")
	var ret = []*model.Milvus{}
	if getCache(cacheKey, &ret) {
		return ret, nil
	}
	ret, err := c.service.ListAll(ctx, opt.Cluster)
	if err != nil {
		err = errors.Wrap(err, "failed to get pods")
		log.Error(err)
		return nil, err
	}
	addCache(cacheKey, &ret)
	return ret, nil
}

var resultCache = cache.New(time.Second*10, time.Second*10)

func addCache(key string, result interface{}) {
	resultCache.Add(key, result, cache.DefaultExpiration)
}

func getCache(key string, res interface{}) bool {
	result, found := resultCache.Get(key)
	if !found {
		return false
	}
	if reflect.TypeOf(result).AssignableTo(reflect.TypeOf(res)) {
		reflect.ValueOf(res).Elem().Set(reflect.ValueOf(result).Elem())
		return true
	}
	return false
}

func (c *MilvusController) handleDownloadLog(ctx *gin.Context) {
	opt := model.GetMilvusLocateOption(ctx)
	logOpt := model.GetLogOption(ctx)

	var writer io.Writer = ctx.Writer
	readerCloser, err := c.service.DownloadLog(ctx, opt, logOpt)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	defer readerCloser.Close()

	ctx.Status(200)
	ctx.Header("Content-Type", "text/plain")
	ctx.Header("Content-Disposition", `attachment; filename="out.log"`)
	_, err = io.Copy(writer, readerCloser)
	if err != nil && err != io.EOF {
		log.Error(err)
	}
}

func (c *MilvusController) handleGetLogs(ctx *gin.Context) {
	opt := model.GetMilvusLocateOption(ctx)
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 5
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Writer.Header())
	if err != nil {
		err = errors.Wrap(err, "failed to upgrade websocket")
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer ws.Close()
	err = c.service.Logs(ctx, ws, opt)
	if err != nil {
		err = errors.Wrap(err, "failed to get logs")
		log.Error(err)
		return
	}
}

func (c *MilvusController) handleGetPods(ctx *gin.Context) (interface{}, error) {
	opt := model.GetMilvusLocateOption(ctx)
	return c.service.ListPods(ctx, opt)
}

func (c *MilvusController) handleShell(ctx *gin.Context) {
	opt := model.GetMilvusLocateOption(ctx)
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 5
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Writer.Header())
	if err != nil {
		err = errors.Wrap(err, "failed to upgrade websocket")
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer ws.Close()
	err = c.service.Shell(ctx, ws, opt)
	if err != nil {
		err = errors.Wrap(err, "failed to shell")
		log.Error(err)
		return
	}
}
