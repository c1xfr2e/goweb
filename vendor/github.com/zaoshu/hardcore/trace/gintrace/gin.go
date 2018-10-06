package gintrace

import (
	"github.com/gin-gonic/gin"
	"github.com/zaoshu/hardcore/trace"
)

const httpTraceHeader = "X-Request-Id"

// WithRequestID gin middleware to set request id
//
// denyExternal deny X-Request-Id from external http request
func WithRequestID(denyExternal ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get request id from X-Request-Id header
		requestID := c.Request.Header.Get(httpTraceHeader)
		if len(requestID) <= 0 || (len(denyExternal) > 0 && denyExternal[0]) {
			requestID = trace.NewRequestID()
		}
		c.Set(trace.RequestIDKey, requestID)
		c.Header(httpTraceHeader, requestID)
	}
}
