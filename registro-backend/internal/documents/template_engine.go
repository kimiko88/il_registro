package documents

import (
	"strings"
)

type TemplateEngine struct{}

func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

// Render replaces {{variable}} with values from context map
func (te *TemplateEngine) Render(template string, data map[string]string) string {
	result := template
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
