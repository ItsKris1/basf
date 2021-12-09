package handler

import (
	"forum/internal/env"
	"forum/internal/session"
	"forum/internal/tpl"
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

		tpl.RenderTemplates(w, "index.html", homeData, "./templates/base.html", "./templates/index.html")

	}
}
