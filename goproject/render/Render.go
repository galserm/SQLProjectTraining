package render

import(
	"io"
	"github.com/labstack/echo"
	"html/template"
)

type TemplateRenderer struct {
    Templates *template.Template
}

func (t *TemplateRenderer) Render(writer io.Writer, name string, data interface {}, ctx echo.Context) error {
    if viewContext, isMap := data.(map[string]interface{}); isMap {
        viewContext["reverse"] = ctx.Echo().Reverse
    }
	return t.Templates.ExecuteTemplate(writer, name, data)
}