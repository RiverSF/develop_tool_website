package template

import (
	"html/template"
	"io"
	"sync"

	"develop_tools/pkg/logger"
	"develop_tools/pkg/path"
)

type MyTemplate struct {
	*template.Template
}

var parsedTemplates sync.Map // map[string]*MyTemplate

func ParseFiles(file string) (*MyTemplate, error) {
	if cached, ok := parsedTemplates.Load(file); ok {
		return cached.(*MyTemplate), nil
	}

	header := path.Join("web", "templates", "_header.html")
	footer := path.Join("web", "templates", "_footer.html")
	tmpl, err := template.ParseFiles(header, footer, file)
	if err != nil {
		return nil, err
	}
	myTmpl := &MyTemplate{Template: tmpl}

	actual, _ := parsedTemplates.LoadOrStore(file, myTmpl)
	return actual.(*MyTemplate), nil
}

func (t MyTemplate) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	err := t.Template.ExecuteTemplate(wr, name, data)
	if err != nil {
		logger.Error("fail to ExecuteTemplate: %s", err.Error())
	}
	return err
}
