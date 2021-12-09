package handler

import (
	"forum/internal/handler/auth"
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type RegisterPage struct {
	AuthInfo auth.RegisterInformation
	UserInfo session.User
}

func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registerPage := RegisterPage{
			AuthInfo: auth.RegInfo,
			UserInfo: session.UserInfo,
		}

		tpl.RenderTemplates(w, "register.html", registerPage, "./templates/register.html", "./templates/base.html")

		auth.RegInfo = auth.RegisterInformation{} // Reset the login messages or they wont change upon reloading the page
	}

}
