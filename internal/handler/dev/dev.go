package dev

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/base"
	"develop_tools/internal/handler/render"
	"develop_tools/pkg/json2struct"
)

func DevJson2GoStruct(c *gin.Context) {
	render.Page(c, "dev_json2gostruct.html")
}

type json2structRequest struct {
	Json string `json:"json"`
}

func DevJson2GoStructApi(c *gin.Context) {
	var req json2structRequest
	if err := c.ShouldBind(&req); err != nil {
		base.ResponseError(c, err.Error())
		return
	}

	result := json2struct.Json2struct(req.Json)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"result": result,
	})
}
