package ctrl

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/haorenfsa/milvus-ops/errs"
)

func parseIntParam(c *gin.Context, name string) (ret int64, err error) {
	raw := c.Param(name)
	ret, err = strconv.ParseInt(raw, 10, 64)
	if err != nil {
		err = errs.ErrBadRequest
	}
	return
}
