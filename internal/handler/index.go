package handler

import (
	"forum/internal/env"
	"forum/internal/errors"
	"forum/internal/session"
	"html/template"
	"net/http"
)

type HomePage struct {
	UserInfo session.User
}

func Index(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session.Check(env.DB, w, r)

		homeData := HomePage{
			UserInfo: session.UserInfo,
		}

		tpl, err := template.ParseFiles("./templates/base.html", "./templates/index.html")
		errors.Check500(w, err)

		err = tpl.ExecuteTemplate(w, "index.html", homeData)
		errors.Check500(w, err)
	}
}
