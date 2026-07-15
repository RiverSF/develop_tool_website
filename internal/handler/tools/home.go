package tools

import (
	"github.com/gin-gonic/gin"

	"develop_tools/internal/handler/render"
)

func Home(c *gin.Context) {
	render.Page(c, "_home.html")
}
