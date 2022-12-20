package template

import (
	"bytes"
	"text/template"
)

func TemplateExecute(template *template.Template, data map[string]any) (string, error) {
	var buffer bytes.Buffer
	err := template.Execute(&buffer, data)

	if err != nil {
		return "", err
	}

	return buffer.String(), err
}
