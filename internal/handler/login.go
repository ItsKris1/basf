package handler

import (
	"forum/internal/handler/auth"
	"forum/internal/session"
	"html/template"
	"net/http"
)

type LoginPage struct {
	UserInfo  session.User
	LoginAuth auth.LoginInformation
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginPage := LoginPage{
			UserInfo:  session.UserInfo,
			LoginAuth: auth.LoginInfo,
		}

		tpl, err := template.ParseFiles("./templates/login.html", "./templates/base.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = tpl.ExecuteTemplate(w, "login.html", loginPage)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		auth.LoginInfo = auth.LoginInformation{} // Reset the login messages or they wont change upon reloading the

	}

}
