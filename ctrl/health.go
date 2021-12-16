package ctrl

import (
	"github.com/gin-gonic/gin"

	"github.com/haorenfsa/milvus-ops/server"
	"github.com/haorenfsa/milvus-ops/service"
)

// HealthController is controller for Task
// implements server.Controller
type HealthController struct {
	service *service.Health
}

func NewHealthController(service *service.Health) *HealthController {
	return &HealthController{service}
}

// Register registers request handler
func (a *HealthController) Register(root gin.IRouter) {
	g := root.Group("/health")
	g.GET("", server.WrapHandler(a.handleGetHealth))
}

func (a *HealthController) handleGetHealth(c *gin.Context) (interface{}, error) {
	return a.service.Get(), nil
}
