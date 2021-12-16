package ctrl

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haorenfsa/milvus-ops/errs"
)

func writeMsgResponseByError(c *gin.Context, err error) {
	if err != nil {
		log.Print("Err:", err)
		code := errs.ErrorToHTTPCode(err)
		c.JSON(code, errMsgMethodFailed)
		return
	}
	c.JSON(http.StatusOK, msgSuccess)
}

func writeObjectResponseByError(c *gin.Context, ret interface{}, err error) {
	if err != nil {
		log.Print("Err:", err)
		code := errs.ErrorToHTTPCode(err)
		c.JSON(code, errMsgMethodFailed)
		return
	}
	c.JSON(http.StatusOK, ret)
}

const errMsgBadBody = "bad body format"
const errMsgMethodFailed = "do method failed"
const msgSuccess = "success"
