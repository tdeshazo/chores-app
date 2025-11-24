package web

import (
	"encoding/json"
	"html/template"
	"path/filepath"
)

// Load templates and add helper funcs.
func ParseTemplates() (*template.Template, error) {
	funcs := template.FuncMap{
		"toJSON": func(v any) template.JS {
			b, err := json.Marshal(v)
			if err != nil {
				return template.JS("null")
			}
			return template.JS(b)
		},
	}

	return template.New("index.html").Funcs(funcs).ParseFiles(filepath.Join("templates", "index.html"))
}
