package frontend

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	Templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func RegisterFrontend(e *echo.Echo) {
	e.Renderer = &TemplateRenderer{
		Templates: template.Must(template.ParseFiles(
			"templates/top.html",
			"templates/basic.html",
			"templates/digest.html",
			"templates/_header.html",
			"templates/_footer.html",
		)),
	}
	e.Static("/static", "static")
}
