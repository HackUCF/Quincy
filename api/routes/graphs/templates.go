/*
Package graphs contains routes that generate html/js graphs that can be embedded into other pages.
*/
package graphs

import (
	"embed"
	"html/template"
)

//go:embed templates
var templateFS embed.FS

var tmpl = template.Must(template.ParseFS(templateFS, "templates/*.html"))
