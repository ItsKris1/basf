package handlers

import (
	"fmt"
	auth "forum/internal/authentication"
	"forum/internal/sessions"
	"html/template"
	"log"
	"net/http"
)

type RegisterPage struct {
	AuthInfo auth.RegisterInformation
	UserInfo sessions.User
}

func Register(w http.ResponseWriter, r *http.Request) {
	registerPage := RegisterPage{
		AuthInfo: auth.RegInfo,
		UserInfo: sessions.UserInfo,
	}
	tpl, err := template.ParseFiles("./templates/register.html", "./templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.ExecuteTemplate(w, "register.html", registerPage)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return
	}

	auth.RegInfo = auth.RegisterInformation{} // Reset the login messages or they wont change upon reloading the page
}
