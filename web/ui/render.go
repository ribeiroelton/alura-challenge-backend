package ui

import (
	"embed"
	"errors"
	"html/template"
	"io"
	"io/fs"

	"github.com/labstack/echo/v4"
)

//go:embed templates
var templatesFS embed.FS

type Template struct {
	templates map[string]*template.Template
}

func newRender() (*Template, error) {
	templates := make(map[string]*template.Template)

	pages, err := fs.Glob(templatesFS, "templates/*-page.tmpl")
	if err != nil {
		return nil, err
	}

	for _, p := range pages {
		templates[p] = template.Must(template.ParseFS(templatesFS, p, "templates/*-layout.tmpl", "templates/*-partials.tmpl"))
	}

	return &Template{
		templates: templates,
	}, nil

}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	file := "templates/" + name

	if _, ok := t.templates[file]; !ok {
		return errors.New("template not found")
	}

	return t.templates[file].ExecuteTemplate(w, name, data)
}
