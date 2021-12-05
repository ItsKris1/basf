package handlers

import (
	"fmt"
	auth "forum/internal/authentication"
	"html/template"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {

	tpl, _ := template.ParseFiles("./templates/register.html")
	err := tpl.Execute(w, auth.RegMsgs)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}

	auth.RegMsgs = auth.RegisterMessages{} // Reset the login messages or they wont change upon reloading the page
}
