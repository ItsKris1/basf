package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("./templates/base.html", "./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}
}
