package ginware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zaoshu/hardcore/logging"
	"github.com/zaoshu/hardcore/logging/runtime"
	"github.com/zaoshu/hardcore/trace"
)

// Logger gin log middleware
//
// Will catch panic
// DO NOT use gin.Default()
// You need create a new gin.Engine
//
//     r := gin.New()
//     r.Use(ginware.Logger())
//
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logging.NewLoggerWithRequestID(trace.GetRequestID(c))
		c.Set(logging.ContextLoggerKey, logger)

		start := time.Now()
		defer func() {
			latency := uint(time.Since(start) / time.Millisecond)
			logger = logger.WithFields(logging.NewStatsField(
				"http",
				map[string]interface{}{
					"url_path": c.Request.URL.Path,
					"method":   c.Request.Method,
					"latency":  latency,
				},
			))

			if err := recover(); err != nil {
				// panic when handle request
				c.AbortWithStatus(http.StatusInternalServerError)
				stack := runtime.Stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)

				msg := fmt.Sprintf(
					"served request [%s], method %s, status %d, latency %d, ip %s",
					c.Request.URL.RequestURI(), c.Request.Method, c.Writer.Status(), latency, c.ClientIP(),
				)
				msg += fmt.Sprintf("\n panic recovered:\n%s\n%s\n%s", string(httprequest), err, stack)
				logger.Error(msg)
				return
			}

			msg := fmt.Sprintf(
				"served request [%s], method %s, status %d, latency %d, ip %s",
				c.Request.URL.RequestURI(), c.Request.Method, c.Writer.Status(), latency, c.ClientIP(),
			)

			if c.Writer.Status() >= http.StatusInternalServerError {
				// internal server error
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				msg += fmt.Sprintf("\n%s\n%s", c.Errors.String(), string(httprequest))
				logger.Error(msg)

			} else if len(c.Errors) > 0 {

				httprequest, _ := httputil.DumpRequest(c.Request, false)
				msg += fmt.Sprintf("\n%s\n%s", c.Errors.String(), string(httprequest))
				logger.Info(msg)

			} else {
				logger.Info(msg)
			}

		}()
		c.Next()
	}
}
