package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// HTTPServer can run to serve http requests
type HTTPServer struct {
	engine     *gin.Engine
	rootRouter gin.IRouter
}

type HandlerFunc func(c *gin.Context) (interface{}, error)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func WrapResponse(data interface{}, err error) *Response {
	if err != nil {
		return &Response{
			Code:    http.StatusInternalServerError, // TODO use actual error code
			Message: err.Error(),
		}
	}
	return &Response{
		Code: http.StatusOK,
		Data: data,
	}
}

func WrapHandler(handle HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handle(c)
		body := WrapResponse(data, err)
		c.JSON(http.StatusOK, body)
	}
}

// NewHTTPServer builds a New HTTPServer
func NewHTTPServer(middlewares ...gin.HandlerFunc) *HTTPServer {
	engine := gin.Default()
	engine.Use(middlewares...)

	return &HTTPServer{
		engine:     engine,
		rootRouter: engine.Group("/api/v1"),
	}
}

// UseControllers registers given controllers to rootRouter
func (h *HTTPServer) UseControllers(ctrls []Controller) {
	for _, ctrl := range ctrls {
		ctrl.Register(h.rootRouter)
	}
}

// ServeStaticPath serves given static path
func (h *HTTPServer) ServeStaticPath(path string) {
	if path == "" {
		return
	}
	fserver := static.LocalFile(path, false)
	server := http.FileServer(fserver)
	h.engine.Use(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/v1") {
			return
		}
		if strings.HasPrefix(ctx.Request.URL.Path, "/app/") {
			ctx.Request.URL.Path = "/"
			ctx.Request.URL.RawPath = "/"
		}
		server.ServeHTTP(ctx.Writer, ctx.Request)
	})
}

// Controller can register request handler
type Controller interface {
	Register(gin.IRouter)
}

// Run runs the server
func (h *HTTPServer) Run(addr string) error {
	return h.engine.Run(addr)
}
