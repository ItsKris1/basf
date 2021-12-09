package tpl

import (
	"net/http"
	"text/template"
)

func RenderTemplates(w http.ResponseWriter, execTpl string, data interface{}, names ...string) {
	tpl, err := template.ParseFiles(names...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = tpl.ExecuteTemplate(w, execTpl, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}
