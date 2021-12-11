package auth

import (
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type LoginPage struct {
	UserInfo  session.User
	LoginAuth LoginInformation
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginPage := LoginPage{
			UserInfo:  session.UserInfo,
			LoginAuth: LoginInfo,
		}

		tpl.RenderTemplates(w, "login.html", loginPage, "./templates/login.html", "./templates/base.html")

		LoginInfo = LoginInformation{} // Reset the login messages or they wont change upon reloading the

	}

}
