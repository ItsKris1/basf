package handlers

import (
	"fmt"
	auth "forum/internal/authentication"
	"forum/internal/sessions"
	"html/template"
	"log"
	"net/http"
)

type LoginPage struct {
	UserInfo  sessions.User
	LoginAuth auth.LoginInformation
}

func Login(w http.ResponseWriter, r *http.Request) {
	loginPage := LoginPage{
		UserInfo:  sessions.UserInfo,
		LoginAuth: auth.LoginInfo,
	}

	tpl, err := template.ParseFiles("./templates/login.html", "./templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.ExecuteTemplate(w, "login.html", loginPage)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500 Internal Server error", 500)
		return

	}

	auth.LoginInfo = auth.LoginInformation{} // Reset the login messages or they wont change upon reloading the page */
}
