package model

import (
	"bytes"
	"os"
	"text/template"
)

type Template struct {
	Name string
	tmpl *template.Template
}

func NewTemplate(filename string) (*Template, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("alert").Parse(string(data))
	if err != nil {
		return nil, err
	}
	return &Template{tmpl: tmpl}, nil
}

func (t *Template) Execute(data interface{}) (string, error) {
	var buf bytes.Buffer
	err := t.tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
