package template

import (
	"fmt"
	"html/template"
	"net/http"
)

// templates to serve
var tpls template.Template

func init() {
	tpls = *template.Must(template.ParseGlob("template/*gohtml"))
}

// ServeTemplate serves a template to the running server
func ServeTemplate(template string, data interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := tpls.ExecuteTemplate(w, template, data); err != nil {
			fmt.Println(err)
		}
	})
}
