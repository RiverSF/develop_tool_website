package template

import (
	"html/template"
	"io"

	"develop_tools/pkg/logger"
	"develop_tools/pkg/path"
)

type MyTemplate struct {
	*template.Template
}

func ParseFiles(file string) (*MyTemplate, error) {
	header := path.Join("web", "templates", "_header.html")
	footer := path.Join("web", "templates", "_footer.html")
	tmpl, err := template.ParseFiles(header, footer, file)
	if err != nil {
		return nil, err
	}
	return &MyTemplate{Template: tmpl}, nil
}

func (t MyTemplate) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	err := t.Template.ExecuteTemplate(wr, name, data)
	if err != nil {
		logger.Error("fail to ExecuteTemplate: %s", err.Error())
	}
	return err
}
