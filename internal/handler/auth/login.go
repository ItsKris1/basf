package auth

import (
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type LoginPage struct {
	UserInfo  session.User
	LoginMsgs LoginMessages
}

func Login(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginPage := LoginPage{
			UserInfo:  session.UserInfo,
			LoginMsgs: LoginMsgs,
		}

		tpl.RenderTemplates(w, "login.html", loginPage, "./templates/auth/login.html", "./templates/base.html")

		LoginMsgs = LoginMessages{} // Reset the login messages or they wont change upon reloading the

	}

}
