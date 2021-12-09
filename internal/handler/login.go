package handler

import (
	"forum/internal/handler/auth"
	"forum/internal/session"
	"forum/internal/tpl"
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

		tpl.RenderTemplates(w, "login.html", loginPage, "./templates/login.html", "./templates/base.html")

		auth.LoginInfo = auth.LoginInformation{} // Reset the login messages or they wont change upon reloading the

	}

}
