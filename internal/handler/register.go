package handler

import (
	"forum/internal/errors"
	"forum/internal/handler/auth"
	"forum/internal/session"
	"html/template"
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

		tpl, err := template.ParseFiles("./templates/register.html", "./templates/base.html")
		errors.Check500(w, err)

		err = tpl.ExecuteTemplate(w, "register.html", registerPage)
		errors.Check500(w, err)

		auth.RegInfo = auth.RegisterInformation{} // Reset the login messages or they wont change upon reloading the page
	}

}
