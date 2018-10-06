package handler

import (
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
)

// Ping handler for testing.
func Ping(c *gin.Context) {
	render.OK(c, nil)
}
