package render

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"develop_tools/pkg/logger"
	"develop_tools/pkg/path"
	"develop_tools/pkg/template"
)

func Page(c *gin.Context, tmplName string) {
	file := path.Join("web", "templates", tmplName)
	t, err := template.ParseFiles(file)
	if err != nil {
		logger.Error("fail to parse template %s: %s", tmplName, err.Error())
		c.String(http.StatusInternalServerError, "template error")
		return
	}
	if err := t.ExecuteTemplate(c.Writer, tmplName, nil); err != nil {
		c.String(http.StatusInternalServerError, "template error")
	}
}
