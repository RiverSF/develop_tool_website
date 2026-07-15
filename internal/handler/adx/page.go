package adx

import (
	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/render"
)

func AdxIndex(c *gin.Context) {
	render.Page(c, "tp_adx.html")
}

func AdxCn(c *gin.Context) {
	render.Page(c, "tp_adx_cn.html")
}

func AdxDSP(c *gin.Context) {
	render.Page(c, "tp_adx_dsp.html")
}

func AdxAdm(c *gin.Context) {
	render.Page(c, "tp_adx_adm.html")
}
