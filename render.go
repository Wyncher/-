package main

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

// struct for custom template
type Template struct {
	templates *template.Template
}

// custom render func
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
