package render

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bluecover/lm/server/errors"
	"github.com/gin-gonic/gin"
	"github.com/zaoshu/hardcore/logging"
)

func OK(c *gin.Context, v interface{}) {
	render(c, http.StatusOK, gin.H{
		"code": 0,
		"data": v,
	})
}

func Fail(c *gin.Context, err error, disableEncode ...bool) {
	c.Error(err)
	if e, ok := err.(*errors.Error); ok {
		render(c, e.HTTPStatus, gin.H{
			"code": e.Code,
			"msg":  e.Message,
		}, disableEncode...)
	} else {
		render(c, http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		}, disableEncode...)
	}
}

func BadRequest(c *gin.Context, err error) {
	c.Error(err)
	render(c, http.StatusBadRequest, gin.H{
		"code": -2,
		"msg":  "Bad Request",
	})
}

func render(c *gin.Context, status int, v interface{}, disableEncode ...bool) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Errorf("json.Marshal %v failed when render, %v", v, err))
	}

	if printData {
		logging.FromContext(c).Infof("[Response][%d] %s", status, string(b))
	}

	if (len(disableEncode) > 0 && disableEncode[0]) || encoder == nil {
		c.Data(status, gin.MIMEJSON, b)
	} else if encoder != nil {
		b, err = encoder.Encode(b)
		if err != nil {
			panic(fmt.Errorf("encode failed when render, %v", err))
		}
		c.Data(status, gin.MIMEPlain, b)
	} else {
		panic("should not go to here")
	}
}
