package auth

import (
	"forum/internal/session"
	"forum/internal/tpl"
	"net/http"
)

type RegisterPage struct {
	AuthInfo RegisterInformation
	UserInfo session.User
}

func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registerPage := RegisterPage{
			AuthInfo: RegInfo,
			UserInfo: session.UserInfo,
		}

		tpl.RenderTemplates(w, "register.html", registerPage, "./templates/register.html", "./templates/base.html")

		RegInfo = RegisterInformation{} // Reset the login messages or they wont change upon reloading the page
	}

}
