package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("./templates/login.html")
	err := tpl.Execute(w, RegValidation) // use RegValidation Success value to see whether registration was succesful;
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}
}
