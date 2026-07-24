package documents

import (
	"html"
	"strings"
)

type TemplateEngine struct{}

func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

// Render replaces {{variable}} with HTML-escaped values from context map to prevent HTML/PDF injection.
func (te *TemplateEngine) Render(template string, data map[string]string) string {
	result := template
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		escaped := html.EscapeString(value)
		result = strings.ReplaceAll(result, placeholder, escaped)
	}
	return result
}
