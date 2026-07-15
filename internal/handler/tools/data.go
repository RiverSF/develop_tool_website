package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/render"
	"develop_tools/pkg/common"
)

func DataDiff(c *gin.Context) {
	render.Page(c, "data_diff.html")
}

func DataEditor(c *gin.Context) {
	render.Page(c, "data_editor.html")
}

func DataFilter(c *gin.Context) {
	render.Page(c, "data_filter.html")
}

func DataCombine(c *gin.Context) {
	render.Page(c, "data_combine.html")
}

func DataCalculator(c *gin.Context) {
	render.Page(c, "data_calculator.html")
}

type Calc struct {
	Command string `json:"command"`
}

func DataCalc(c *gin.Context) {
	var ca Calc
	if err := c.ShouldBind(&ca); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "result": err.Error()})
		return
	}

	result, err := common.EvaluateCalcExpression(ca.Command)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "result": result})
}

func DataUnique(c *gin.Context) {
	render.Page(c, "data_unique.html")
}

type DataUniqueStruct struct {
	Content string `json:"content"`
}

func DataUniqueExec(c *gin.Context) {
	var uq DataUniqueStruct
	if err := c.ShouldBind(&uq); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "result": err.Error()})
		return
	}

	result := common.UniqueCountLines(uq.Content)
	c.JSON(http.StatusOK, gin.H{"status": 200, "result": result})
}

func DataInterval(c *gin.Context) {
	render.Page(c, "data_interval.html")
}

func DataLine(c *gin.Context) {
	render.Page(c, "data_line.html")
}
