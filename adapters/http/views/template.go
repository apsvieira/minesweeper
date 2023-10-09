package views

import (
	"embed"
	"fmt"
	"html/template"
)

var (
	//go:embed "templates/*"
	todoTemplates embed.FS
)

func NewTemplates() (*template.Template, error) {
	t := template.New("").Funcs(template.FuncMap{
		"mapElements": mapElements,
	})

	return t.ParseFS(todoTemplates, "templates/*.html")
}

func mapElements(pairs ...any) (map[string]any, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("mapElements: odd number of elements")
	}

	m := make(map[string]any, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		k, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("mapElements: non-string key: %T", pairs[i])
		}
		m[k] = pairs[i+1]
	}
	return m, nil
}
