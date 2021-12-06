package handlers

import (
	"fmt"
	auth "forum/internal/authentication"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("./templates/index.html")

	err := tpl.Execute(w, auth.LoginInfo.LoggedUser)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}
}
