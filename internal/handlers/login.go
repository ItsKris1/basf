package handlers

import (
	"fmt"
	auth "forum/internal/authentication"
	"html/template"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("./templates/login.html")
	err := tpl.Execute(w, auth.LoginInfo)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return

	}

	auth.LoginInfo = auth.LoginInformation{} // Reset the login messages or they wont change upon reloading the page */
}
