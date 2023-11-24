package http

import (
	"github.com/gin-gonic/gin"
)

func New()*gin.Engine {
	r := gin.Default()
	r.POST("/n", NormalAnalysis)
	return r
}
