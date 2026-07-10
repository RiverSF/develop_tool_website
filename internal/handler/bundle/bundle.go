package bundle

import (
	"strings"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/base"
	"develop_tools/internal/handler/render"
	pkgbundle "develop_tools/pkg/bundle"
)

func BundleIndex(c *gin.Context) {
	render.Page(c, "tp_bundle.html")
}

type extractRequest struct {
	Bundles string `json:"bundles"`
	Fields  string `json:"fields"`
}

func Extract(c *gin.Context) {
	var req extractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		base.ResponseError(c, err.Error())
		return
	}
	if strings.TrimSpace(req.Bundles) == "" {
		base.ResponseError(c, "请输入 bundle")
		return
	}

	results := pkgbundle.Extract(req.Bundles, req.Fields)
	if len(results) == 0 {
		base.ResponseError(c, "未解析到有效 bundle")
		return
	}
	base.ResponseOk(c, results)
}
