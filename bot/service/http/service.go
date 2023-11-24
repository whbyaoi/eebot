package http

import (
	"eebot/bot/service/analysis300"

	"github.com/gin-gonic/gin"
)

type result struct {
	Data interface{}
	Err  string
}

func NormalAnalysis(c *gin.Context) {
	name := c.Query("name")
	rs, err := analysis300.ExportWinOrLoseAnalysisAdvanced(name)
	if err != nil {
		c.JSON(400, result{nil, err.Error()})
		return
	}
	c.JSON(200, result{rs, ""})
}
