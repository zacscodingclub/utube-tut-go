package utils

import (
	"html/template"
	"net/http"
)

var templates *template.Template

func LoadTemplates(p string) {
	templates = template.Must(template.ParseGlob(p))
}

func ExecuteTemplate(w http.ResponseWriter, t string, d interface{}) {
	templates.ExecuteTemplate(w, t, d)
}
